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
	s.router.HandleFunc("/subs/last", s.lastSubView())
	s.router.HandleFunc("/api/subs/last", s.fetchLastSub())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)

}
