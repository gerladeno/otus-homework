package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	internalhttp "github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/server/http"
)

type CalendarHttpApi struct {
	ConnHTTP *http.Client
	Host     string
}

func (a *CalendarHttpApi) CreateEvent(event common.Event) (int64, int) {
	w := bytes.Buffer{}
	err := json.NewEncoder(&w).Encode(event)
	if err != nil {
		return 0, http.StatusInternalServerError
	}
	r, err := a.ConnHTTP.Post(a.Host+"/api/v1/addEvent", "application/json", &w)
	if err != nil {
		return 0, http.StatusInternalServerError
	}
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

func (a *CalendarHttpApi) DeleteEvent(id int64) int {
	r, err := a.ConnHTTP.Get(a.Host + fmt.Sprintf("/api/v1/deleteEvent/%d", id))
	if err != nil {
		return http.StatusInternalServerError
	}
	var result internalhttp.JSONResponse
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return http.StatusInternalServerError
	}
	return result.Code
}

func (a *CalendarHttpApi) UpdateEvent(event common.Event, id int64) int {
	w := bytes.Buffer{}
	err := json.NewEncoder(&w).Encode(event)
	if err != nil {
		return http.StatusInternalServerError
	}
	r, err := a.ConnHTTP.Post(a.Host+fmt.Sprintf("/api/v1/editEvent/%d", id), "application/json", &w)
	if err != nil {
		return http.StatusInternalServerError
	}
	var result internalhttp.JSONResponse
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return http.StatusInternalServerError
	}
	return result.Code
}

func (a *CalendarHttpApi) ListEventsByDay(date string) ([]common.Event, int) {
	r, err := a.ConnHTTP.Get(a.Host + fmt.Sprintf("/api/v1/listEventsByDay?date=%s", date))
	if err != nil {
		return nil, http.StatusInternalServerError
	}
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

func (a *CalendarHttpApi) ListEventsByWeek(date string) ([]common.Event, int) {
	r, err := a.ConnHTTP.Get(a.Host + fmt.Sprintf("/api/v1/listEventsByWeek?date=%s", date))
	if err != nil {
		return nil, http.StatusInternalServerError
	}
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

func (a *CalendarHttpApi) ListEventsByMonth(date string) ([]common.Event, int) {
	r, err := a.ConnHTTP.Get(a.Host + fmt.Sprintf("/api/v1/listEventsByMonth?date=%s", date))
	if err != nil {
		return nil, http.StatusInternalServerError
	}
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
