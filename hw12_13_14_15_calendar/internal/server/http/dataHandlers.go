package internalhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/go-chi/chi"
)

const dateFormat = "2006-01-02"

var (
	ErrWrongEventID     = errors.New("invalid or empty id")
	ErrEmptyRequestBody = errors.New("empty request body")
	ErrUnparsableEvent  = errors.New("err parsing event")
)

type JSONResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *string     `json:"error,omitempty"`
	Code  int         `json:"code"`
}

type ID struct {
	ID int64 `json:"id"`
}

func (h *EventHandler) listEventsByDayHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		writeErrResponse(w, "unparsable date, use YYYY-MM-DD: "+dateFormat, http.StatusBadRequest)
	}
	events, err := h.app.ListEventsByDay(r.Context(), date)
	if err != nil {
		h.log.Warn("failed to get list of events: ", err)
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, events)
}

func (h *EventHandler) listEventsByWeekHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		writeErrResponse(w, "unparsable date, use YYYY-MM-DD: "+dateFormat, http.StatusBadRequest)
	}
	events, err := h.app.ListEventsByWeek(r.Context(), date)
	if err != nil {
		h.log.Warn("failed to get list of events: ", err)
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, events)
}

func (h *EventHandler) listEventsByMonthHandler(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		writeErrResponse(w, "unparsable date, use YYYY-MM-DD: "+dateFormat, http.StatusBadRequest)
	}
	events, err := h.app.ListEventsByMonth(r.Context(), date)
	if err != nil {
		h.log.Warn("failed to get list of events: ", err)
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, events)
}

func (h *EventHandler) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeErrResponse(w, ErrWrongEventID.Error(), http.StatusBadRequest)
		h.log.Debug(err)
		return
	}
	err = h.app.DeleteEvent(r.Context(), id)
	if err != nil {
		if errors.Is(err, common.ErrNoSuchEvent) {
			h.log.Debugf("failed to remove an event %d: %s", id, err.Error())
			writeErrResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		h.log.Warnf("failed to remove an event %d: %s", id, err.Error())
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, ID{ID: id})
}

func (h *EventHandler) addEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		h.log.Debug("empty request body")
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}
	event := new(common.Event)
	if err := event.ParseEvent(r); err != nil {
		h.log.Debug("can't parse events: ", err)
		writeErrResponse(w, ErrUnparsableEvent.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.app.CreateEvent(r.Context(), event)
	if err != nil {
		h.log.Warn("failed to add event: ", err)
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, ID{ID: id})
}

func (h *EventHandler) editEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeErrResponse(w, ErrWrongEventID.Error(), http.StatusBadRequest)
		h.log.Debug(err)
		return
	}
	if r.Body == nil {
		h.log.Debug("empty request body")
		writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
		return
	}
	event := new(common.Event)
	if err := event.ParseEvent(r); err != nil {
		h.log.Debug("can't parse events: ", err)
		writeErrResponse(w, ErrUnparsableEvent.Error(), http.StatusBadRequest)
		return
	}
	err = h.app.UpdateEvent(r.Context(), id, event)
	if err != nil {
		if errors.Is(err, common.ErrNoSuchEvent) {
			h.log.Debugf("failed to edit an event %d: %s", id, err.Error())
			writeErrResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		h.log.Warnf("failed to edit an event %d: %s", id, err.Error())
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, ID{ID: id})
}

func writeOkResponse(w http.ResponseWriter, data interface{}) {
	status := http.StatusOK
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	response := JSONResponse{
		Data: &data,
		Code: status,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func writeErrResponse(w http.ResponseWriter, err string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	response := JSONResponse{
		Error: &err,
		Code:  status,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func parseIDParam(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, ErrWrongEventID
	}
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
