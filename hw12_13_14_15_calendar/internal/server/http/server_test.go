package internalhttp

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type TestApp struct{}

func (t TestApp) ReadEvent(_ context.Context, id uint64) (event *common.Event, err error) {
	switch id {
	case 0:
		err = common.ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return event, err
}

func (t TestApp) CreateEvent(_ context.Context, event *common.Event) (uint64, error) {
	if event.ID == 0 {
		return 1, nil
	}
	return 0, io.ErrShortBuffer
}

func (t TestApp) UpdateEvent(_ context.Context, _ *common.Event, id uint64) (err error) {
	switch id {
	case 0:
		err = common.ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return err
}

func (t TestApp) DeleteEvent(_ context.Context, id uint64) (err error) {
	switch id {
	case 0:
		err = common.ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return err
}

func (t TestApp) ListEvents(_ context.Context) ([]*common.Event, error) {
	return make([]*common.Event, 5), nil
}

func TestHandlers(t *testing.T) {
	ts := NewServer(TestApp{}, logrus.New(), "test_version", 3001)
	t.Run("listEntries", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", "/api/v1/listEvents", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ts.listEventsHandler)
		handler.ServeHTTP(rr, req)
		require.Equal(t, rr.Code, http.StatusOK)
	})
}
