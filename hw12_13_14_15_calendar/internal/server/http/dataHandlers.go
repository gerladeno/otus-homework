package internalhttp

import (
	"encoding/json"
	"errors"
	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var ErrWrongEventId = errors.New("invalid or empty id")
var ErrEmptyRequestBody = errors.New("empty request body")
var ErrUnparsableEvent = errors.New("err parsing event")

type JsonResponse struct {
	Data  *interface{} `json:"data,omitempty"`
	Error *string     `json:"error,omitempty"`
	Code  int         `json:"code"`
}

type Id struct {
	Id uint64 `json:"id"`
}

func listEventsHandler(storage common.Storage, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		events, err := storage.ListEvents()
		if err != nil {
			log.Warn("failed to get list of events: ", err)
			writeErrResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeOkResponse(w, events, http.StatusOK)
	}
}

func getEventHandler(storage common.Storage, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIdParam(r)
		if err != nil {
			writeErrResponse(w, ErrWrongEventId.Error(), http.StatusBadRequest)
			log.Debug(err)
			return
		}
		event, err := storage.GetEvent(id)
		if err != nil {
			if errors.Is(err, common.ErrNoSuchEvent) {
				log.Debug("failed to get an event %d: %s", id, err.Error())
				writeErrResponse(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Warnf("failed to get an event %d: %s", id, err.Error())
			writeErrResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeOkResponse(w, event, http.StatusOK)
	}
}

func removeEventHandler(storage common.Storage, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIdParam(r)
		if err != nil {
			writeErrResponse(w, ErrWrongEventId.Error(), http.StatusBadRequest)
			log.Debug(err)
			return
		}
		err = storage.RemoveEvent(r.Context(), id)
		if err != nil {
			if errors.Is(err, common.ErrNoSuchEvent) {
				log.Debug("failed to remove an event %d: %s", id, err.Error())
				writeErrResponse(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Warnf("failed to remove an event %d: %s", id, err.Error())
			writeErrResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeOkResponse(w, Id{Id: id}, http.StatusOK)
	}
}

func addEventHandler(storage common.Storage, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			log.Debug("empty request body")
			writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
			return
		}
		event := new(common.Event)
		if err := event.ParseEvent(r); err != nil {
			log.Debug("can't parse events: ", err)
			writeErrResponse(w, ErrUnparsableEvent.Error(), http.StatusBadRequest)
			return
		}
		id, err := storage.AddEvent(r.Context(), *event)
		if err != nil {
			log.Warn("failed to add event: ", err)
			writeErrResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeOkResponse(w, Id{Id: id}, http.StatusOK)
	}
}

func editEventHandler(storage common.Storage, log *logrus.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIdParam(r)
		if err != nil {
			writeErrResponse(w, ErrWrongEventId.Error(), http.StatusBadRequest)
			log.Debug(err)
			return
		}
		if r.Body == nil {
			log.Debug("empty request body")
			writeErrResponse(w, ErrEmptyRequestBody.Error(), http.StatusBadRequest)
			return
		}
		event := new(common.Event)
		if err := event.ParseEvent(r); err != nil {
			log.Debug("can't parse events: ", err)
			writeErrResponse(w, ErrUnparsableEvent.Error(), http.StatusBadRequest)
			return
		}
		err = storage.EditEvent(r.Context(), id, *event)
		if err != nil {
			if errors.Is(err, common.ErrNoSuchEvent) {
				log.Debug("failed to edit an event %d: %s", id, err.Error())
				writeErrResponse(w, err.Error(), http.StatusNotFound)
				return
			}
			log.Warnf("failed to edit an event %d: %s", id, err.Error())
			writeErrResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeOkResponse(w, Id{Id: id}, http.StatusOK)
	}
}

func writeOkResponse(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	response := JsonResponse{
		Data: &data,
		Code: status,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func writeErrResponse(w http.ResponseWriter, err string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	response := JsonResponse{
		Error: &err,
		Code:  status,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func parseIdParam(r *http.Request) (uint64, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, ErrWrongEventId
	}
	id, err := strconv.ParseUint(idStr, 0, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
