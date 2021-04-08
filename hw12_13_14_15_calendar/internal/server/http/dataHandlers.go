package internalhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/go-chi/chi"
)

var (
	ErrWrongEventID     = errors.New("invalid or empty id")
	ErrEmptyRequestBody = errors.New("empty request body")
	ErrUnparsableEvent  = errors.New("err parsing event")
)

type JSONResponse struct {
	Data  *interface{} `json:"data,omitempty"`
	Error *string      `json:"error,omitempty"`
	Code  int          `json:"code"`
}

type ID struct {
	ID uint64 `json:"id"`
}

func (h *EventHandler) listEventsHandler(w http.ResponseWriter, r *http.Request) {
	events, err := h.app.ListEvents(r.Context())
	if err != nil {
		h.log.Warn("failed to get list of events: ", err)
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, events)
}

func (h *EventHandler) getEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		writeErrResponse(w, ErrWrongEventID.Error(), http.StatusBadRequest)
		h.log.Debug(err)
		return
	}
	event, err := h.app.ReadEvent(r.Context(), id)
	if err != nil {
		if errors.Is(err, common.ErrNoSuchEvent) {
			h.log.Debugf("failed to get an event %d: %s", id, err.Error())
			writeErrResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		h.log.Warnf("failed to get an event %d: %s", id, err.Error())
		writeErrResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeOkResponse(w, event)
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
	err = h.app.UpdateEvent(r.Context(), event, id)
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

func parseIDParam(r *http.Request) (uint64, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, ErrWrongEventID
	}
	id, err := strconv.ParseUint(idStr, 0, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
