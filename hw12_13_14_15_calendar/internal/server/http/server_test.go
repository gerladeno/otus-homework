package internalhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestListHandler(t *testing.T) {
	log := logrus.New()
	th := NewEventHandler(common.TestApp{}, log)
	tr := NewRouter(th, log, "test")
	var result JSONResponse

	t.Run("listEntriesByDay", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/listEventsByDay?date=1987-10-16", nil)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusOK)
		err := json.NewDecoder(w.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, len(result.Data.([]interface{})), 5)
	})
	t.Run("listEntriesByDay err", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/listEventsByDay?date=1987-15-12", nil)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusBadRequest)
	})
	t.Run("listEntriesByWeek", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/listEventsByWeek?date=1987-10-16", nil)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusOK)
		err := json.NewDecoder(w.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, len(result.Data.([]interface{})), 15)
	})
	t.Run("listEntriesByWeek err", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/listEventsByWeek?date=1987-15-12", nil)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusBadRequest)
	})
	t.Run("listEntriesByMonth", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/listEventsByMonth?date=1987-10-16", nil)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusOK)
		err := json.NewDecoder(w.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, len(result.Data.([]interface{})), 50)
	})
	t.Run("listEntriesByMonth err", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/listEventsByMonth?date=1987-15-12", nil)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusBadRequest)
	})
}

func TestDeleteHandler(t *testing.T) {
	log := logrus.New()
	th := NewEventHandler(common.TestApp{}, log)
	tr := NewRouter(th, log, "test")
	var result JSONResponse

	testsDelete := []struct {
		name    string
		id      int
		errCode int
		err     string
	}{
		{"no such entry", 0, http.StatusNotFound, "no such event"},
		{"internal error", 1, http.StatusInternalServerError, "short buffer"},
	}
	for _, test := range testsDelete {
		test := test
		t.Run("deleteEntry", func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/deleteEvent/%d", test.id), nil)
			tr.ServeHTTP(w, r)
			require.Equal(t, w.Code, test.errCode)
			err := json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)
			require.Equal(t, *result.Error, test.err)
		})
	}
	t.Run("ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/v1/deleteEvent/2", nil)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusOK)
	})
}

func TestCreateEvent(t *testing.T) {
	log := logrus.New()
	th := NewEventHandler(common.TestApp{}, log)
	tr := NewRouter(th, log, "test")
	var result struct {
		Data struct {
			ID int `json:"id"`
		} `json:"data,omitempty"`
		Error *string `json:"error,omitempty"`
		Code  int     `json:"code"`
	}
	t.Run("ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := bytes.NewReader([]byte(`{"id":0, "title":"jopa","startTime":"2021-04-08T22:54:10+03:00","duration":300}`))
		r := httptest.NewRequest("POST", "/api/v1/addEvent", body)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusOK)
		err := json.NewDecoder(w.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, result.Data.ID, 1)
	})
	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := bytes.NewReader([]byte(`{"id":1, "title":"opiat jopa"}`))
		r := httptest.NewRequest("POST", "/api/v1/addEvent", body)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusInternalServerError)
	})
}

func TestUpdateEvent(t *testing.T) {
	log := logrus.New()
	th := NewEventHandler(common.TestApp{}, log)
	tr := NewRouter(th, log, "test")
	var result struct {
		Data *struct {
			ID int `json:"id"`
		} `json:"data,omitempty"`
		Error *string `json:"error,omitempty"`
		Code  int     `json:"code"`
	}
	testsRead := []struct {
		name    string
		id      int
		errCode int
		err     string
	}{
		{"no such entry", 0, http.StatusNotFound, "no such event"},
		{"internal error", 1, http.StatusInternalServerError, "short buffer"},
	}
	for _, test := range testsRead {
		test := test
		t.Run("readEntry", func(t *testing.T) {
			w := httptest.NewRecorder()
			body := bytes.NewReader([]byte(`{"id":0, "title":"jopa","startTime":"2021-04-08T22:54:10+03:00","duration":300}`))
			r := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/editEvent/%d", test.id), body)
			tr.ServeHTTP(w, r)
			require.Equal(t, w.Code, test.errCode)
			err := json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)
			require.Equal(t, *result.Error, test.err)
			require.Nil(t, result.Data)
		})
	}
	t.Run("ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := bytes.NewReader([]byte(`{"id":0, "title":"jopa","startTime":"2021-04-08T22:54:10+03:00","duration":300}`))
		r := httptest.NewRequest("POST", "/api/v1/editEvent/2", body)
		tr.ServeHTTP(w, r)
		require.Equal(t, w.Code, http.StatusOK)
		err := json.NewDecoder(w.Body).Decode(&result)
		require.NoError(t, err)
		require.Equal(t, result.Data.ID, 2)
	})
}
