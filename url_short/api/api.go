package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5" // Using Chi because it allows for URL parameters and easier middleware
	"github.com/go-chi/chi/v5/middleware"
	"github.com/its-kos/assignment1/storage"
)

type ApiServer struct {
	addr         string
	db           storage.Storage
	startTime    time.Time
	requestCount int
	succesful    int
}

func NewApiServer(addr string, db storage.Storage) *ApiServer {
	return &ApiServer{addr: addr, db: db, startTime: time.Now(), requestCount: 0, succesful: 0}
}

func (s *ApiServer) Run() error {
	router := chi.NewRouter()
	router.Use(middleware.Logger) // <--<< Logger should come before Recoverer
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/status"))

	router.Get("/", s.handleGetAllIDs)
	router.Get("/{id}", s.handleGetURLByID)
	router.Get("/metrics", s.handleGetMetrics)
	router.Get("/metrics/{id}", s.handleGetMetricsByID)
	router.Post("/", s.handleCreateURLAlias)
	router.Put("/{id}", s.handleUpdateURLByID)
	router.Delete("/{id}", s.handleDeleteURLByID)
	router.Delete("/", s.handleDeleteAllURLs)

	server := &http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	log.Printf("Server is listening on %s", s.addr)

	return server.ListenAndServe()
}
