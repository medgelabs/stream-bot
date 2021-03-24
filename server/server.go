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
	debugClient        *DebugClient
}

// New returns a Server instance to be run with http.ListenAndServe()
func New(metricStore cache.Cache, alertWebSocket *bot.WriteOnlyUnsafeWebSocket, debugClient *DebugClient) *Server {
	srv := &Server{
		router:             http.NewServeMux(),
		viewerMetricsStore: metricStore,
		alertWebSocket:     alertWebSocket,
		debugClient:        debugClient,
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

	s.router.HandleFunc("/api/bits/last", s.fetchLastBits())
	s.router.HandleFunc("/bits/last", s.lastBitsView(baseURL+"/api/bits/last"))

	// DEBUG - trigger various events for testing
	// TODO how do secure when deploy?
	s.router.HandleFunc("/debug/sub", s.debugSub(s.debugClient))
	s.router.HandleFunc("/debug/bit", s.debugBit(s.debugClient))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
