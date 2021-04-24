package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	internalhttp "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/http"
)

type CalendarHTTPApi struct {
	ConnHTTP *http.Client
	Host     string
}

func (a *CalendarHTTPApi) CreateEvent(ctx context.Context, event common.Event) (int64, int) {
	w := bytes.Buffer{}
	err := json.NewEncoder(&w).Encode(event)
	if err != nil {
		return 0, http.StatusInternalServerError
	}
	req, err := http.NewRequestWithContext(ctx, "POST", a.Host+"/api/v1/addEvent", &w)
	if err != nil {
		return 0, http.StatusInternalServerError
	}
	r, err := a.ConnHTTP.Do(req)
	if err != nil {
		return 0, http.StatusInternalServerError
	}
	defer func() {
		_ = r.Body.Close()
	}()
	var result struct {
		Data struct {
			ID int64 `json:"id"`
		} `json:"data,omitempty"`
		Error *string `json:"error,omitempty"`
		Code  int     `json:"code"`
	}
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return 0, http.StatusInternalServerError
	}
	return result.Data.ID, result.Code
}

func (a *CalendarHTTPApi) DeleteEvent(ctx context.Context, id int64) int {
	req, err := http.NewRequestWithContext(ctx, "GET", a.Host+fmt.Sprintf("/api/v1/deleteEvent/%d", id), nil)
	if err != nil {
		return http.StatusInternalServerError
	}
	r, err := a.ConnHTTP.Do(req)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer func() {
		_ = r.Body.Close()
	}()
	var result internalhttp.JSONResponse
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return http.StatusInternalServerError
	}
	return result.Code
}

func (a *CalendarHTTPApi) UpdateEvent(ctx context.Context, event common.Event, id int64) int {
	w := bytes.Buffer{}
	err := json.NewEncoder(&w).Encode(event)
	if err != nil {
		return http.StatusInternalServerError
	}
	req, err := http.NewRequestWithContext(ctx, "POST", a.Host+fmt.Sprintf("/api/v1/editEvent/%d", id), &w)
	if err != nil {
		return http.StatusInternalServerError
	}
	r, err := a.ConnHTTP.Do(req)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer func() {
		_ = r.Body.Close()
	}()
	var result internalhttp.JSONResponse
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return http.StatusInternalServerError
	}
	return result.Code
}

func (a *CalendarHTTPApi) ListEventsByDay(ctx context.Context, date string) ([]common.Event, int) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.Host+fmt.Sprintf("/api/v1/listEventsByDay?date=%s", date), nil)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	r, err := a.ConnHTTP.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer func() {
		_ = r.Body.Close()
	}()
	var result struct {
		Data  []common.Event `json:"data,omitempty"`
		Error *string        `json:"error,omitempty"`
		Code  int            `json:"code"`
	}
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return result.Data, result.Code
}

func (a *CalendarHTTPApi) ListEventsByWeek(ctx context.Context, date string) ([]common.Event, int) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.Host+fmt.Sprintf("/api/v1/listEventsByWeek?date=%s", date), nil)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	r, err := a.ConnHTTP.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer func() {
		_ = r.Body.Close()
	}()
	var result struct {
		Data  []common.Event `json:"data,omitempty"`
		Error *string        `json:"error,omitempty"`
		Code  int            `json:"code"`
	}
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return result.Data, result.Code
}

func (a *CalendarHTTPApi) ListEventsByMonth(ctx context.Context, date string) ([]common.Event, int) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.Host+fmt.Sprintf("/api/v1/listEventsByMonth?date=%s", date), nil)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	r, err := a.ConnHTTP.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer func() {
		_ = r.Body.Close()
	}()
	var result struct {
		Data  []common.Event `json:"data,omitempty"`
		Error *string        `json:"error,omitempty"`
		Code  int            `json:"code"`
	}
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return result.Data, result.Code
}
