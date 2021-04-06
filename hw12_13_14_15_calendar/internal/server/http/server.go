package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gerladeno/otus_homeworks/hw12_13_14_15_calendar/internal/storage/common"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Server struct {
	app     Application
	storage common.Storage
	log     *logrus.Logger
	router  chi.Router
}

type Application interface { // TODO
}

func NewServer(app Application, storage common.Storage, log *logrus.Logger, version interface{}) *Server {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(loggingMiddleware(log))
	r.NotFound(notFoundHandler)
	r.Get("/hello", helloHandler)
	r.Get("/version", versionHandler(version))
	return &Server{
		app:     app,
		storage: storage,
		log:     log,
		router:  r,
	}
}

func (s *Server) Start(ctx context.Context) error {
	port := ":3000"
	server := &http.Server{
		Addr:              port,
		Handler:           s.router,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		ctx.Done()
		return err
	}
	s.log.Infof("started server on %s", port)
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	ctx.Done()
	return nil
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
