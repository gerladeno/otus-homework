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

type TestStorage struct{}

func (t TestStorage) ReadEvent(_ context.Context, id uint64) (event *common.Event, err error) {
	switch id {
	case 0:
		err = common.ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return event, err
}

func (t TestStorage) CreateEvent(_ context.Context, event common.Event) (uint64, error) {
	if event.ID == 0 {
		return 1, nil
	}
	return 0, io.ErrShortBuffer
}

func (t TestStorage) UpdateEvent(_ context.Context, id uint64, _ common.Event) (err error) {
	switch id {
	case 0:
		err = common.ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return err
}

func (t TestStorage) DeleteEvent(_ context.Context, id uint64) (err error) {
	switch id {
	case 0:
		err = common.ErrNoSuchEvent
	case 1:
		err = io.ErrShortBuffer
	default:
	}
	return err
}

func (t TestStorage) ListEvents(_ context.Context) ([]common.Event, error) {
	return make([]common.Event, 5), nil
}

func TestHandlers(t *testing.T) {
	var ts TestStorage
	t.Run("listEntries", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", "/api/v1/listEvents", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(listEventsHandler(ts, logrus.New()))
		handler.ServeHTTP(rr, req)
		require.Equal(t, rr.Code, http.StatusOK)
	})
}
