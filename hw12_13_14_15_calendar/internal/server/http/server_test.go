package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestListHandler(t *testing.T) {
	ts := NewServer(TestApp{}, logrus.New(), "test_version", 3001)
	w := httptest.NewRecorder()
	var result JSONResponse

	t.Run("listEntries", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/api/v1/listEvents", nil)
		ts.listEventsHandler(w, r)
		require.Equal(t, w.Code, http.StatusOK)
		err := json.NewDecoder(w.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, len((*result.Data).([]interface{})), 5)
	})
}

func TestDeleteHandler(t *testing.T) {
	ts := NewServer(TestApp{}, logrus.New(), "test_version", 3001)
	w := httptest.NewRecorder()
	var result JSONResponse

	testsDelete := []struct {
		name    string
		id      int
		errCode int
		err     string
	}{
		//{"invalid id",-1, http.StatusBadRequest, "invalid or empty id"},
		{"no such entry", 0, http.StatusNotFound, ""},
		//{"internal error",1, http.StatusInternalServerError},
		//{"ok",2, http.StatusOK},
	}
	for _, test := range testsDelete {
		test := test
		t.Run("deleteEntry", func(t *testing.T) {
			r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/deleteEvent/%d", test.id), nil)
			ts.deleteEventHandler(w, r)
			require.Equal(t, w.Code, test.errCode)
			err := json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)
			require.Equal(t, *result.Error, test.err)
		})
	}
}

//body := bytes.NewReader([]byte(fmt.Sprintf(`{"id":%d}`, test.id)))
