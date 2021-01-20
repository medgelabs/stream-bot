package server

import (
	"net/http"
)

type server struct {
	router *http.ServeMux
}

func New() *server {
	srv := &server{
		router: http.NewServeMux(),
	}

	srv.routes()
	return srv
}

func (s *server) routes() {
	// TODO pull from config
	baseUrl := "http://localhost:8080"

	s.router.HandleFunc("/api/subs/last", s.fetchLastSub())
	s.router.HandleFunc("/subs/last", s.lastSubView(baseUrl+"/api/subs/last"))
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
