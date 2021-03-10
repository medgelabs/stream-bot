package server

import (
	"medgebot/bot"
	"medgebot/cache"
	"net/http"
)

// Server REST API
type Server struct {
	router             *http.ServeMux
	viewerMetricsStore cache.Cache
	alertWebSocket     *bot.WriteOnlyUnsafeWebSocket
}

// New returns a Server instance to be run with http.ListenAndServe()
func New(metricStore cache.Cache, alertWebSocket *bot.WriteOnlyUnsafeWebSocket) *Server {
	srv := &Server{
		router:             http.NewServeMux(),
		viewerMetricsStore: metricStore,
		alertWebSocket:     alertWebSocket,
	}

	// TODO receive WebSocket connection for alerts and assign to alertWebSocket

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
