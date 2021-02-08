package server

import (
	"medgebot/bot/viewer"
	"net/http"
)

// Server REST API
type Server struct {
	router             *http.ServeMux
	viewerMetricsStore viewer.MetricStore
}

// New returns a Server instance to be run with http.ListenAndServe()
func New(metricStore viewer.MetricStore) *Server {
	srv := &Server{
		router:             http.NewServeMux(),
		viewerMetricsStore: metricStore,
	}

	srv.routes()
	return srv
}

func (s *Server) routes() {
	// TODO pull from config
	baseURL := "http://localhost:8080"

	s.router.HandleFunc("/api/subs/last", s.fetchLastSub())
	s.router.HandleFunc("/subs/last", s.lastSubView(baseURL+"/api/subs/last"))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
