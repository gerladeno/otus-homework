package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/common"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router chi.Router
	port   string
	log    *logrus.Logger
	server *http.Server
}

type EventHandler struct {
	app common.Application
	log *logrus.Logger
}

func NewEventHandler(app common.Application, log *logrus.Logger) *EventHandler {
	return &EventHandler{app: app, log: log}
}

func NewRouter(handler *EventHandler, log *logrus.Logger, version interface{}) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(15 * time.Second))
	r.NotFound(notFoundHandler)
	r.Get("/hello", helloHandler)
	r.Get("/version", versionHandler(version))
	r.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(loggingMiddleware(log))
			r.Route("/v1", func(r chi.Router) {
				r.Get("/listEventsByDay", handler.listEventsByDayHandler)
				r.Get("/listEventsByWeek", handler.listEventsByWeekHandler)
				r.Get("/listEventsByMonth", handler.listEventsByMonthHandler)
				r.Get("/deleteEvent/{id}", handler.deleteEventHandler)
				r.Post("/addEvent", handler.addEventHandler)
				r.Post("/editEvent/{id}", handler.editEventHandler)
			})
		})
	})
	return r
}

func NewServer(log *logrus.Logger, r chi.Router, port int) *Server {
	server := Server{
		router: r,
		port:   ":" + strconv.Itoa(port),
		log:    log,
	}
	return &server
}

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{
		Addr:              s.port,
		Handler:           s.router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
	}
	go func() {
		<-ctx.Done()
		if err := s.Stop(); err != nil {
			s.log.Warnf("failed to stop http server: %s", err.Error())
		}
	}()
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	return s.server.Close()
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "404 page not found,", http.StatusNotFound)
}

func versionHandler(version interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(version)
	}
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode("Hello, world!")
}
