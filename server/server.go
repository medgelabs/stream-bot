package server

import (
	"medgebot/bot"
	"medgebot/cache"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server REST API
type Server struct {
	router             *chi.Mux
	viewerMetricsStore cache.Cache
	alertWebSocket     *bot.WriteOnlyUnsafeWebSocket
	debugClient        *DebugClient
}

// New returns a Server instance to be run with http.ListenAndServe()
func New(metricStore cache.Cache, alertWebSocket *bot.WriteOnlyUnsafeWebSocket, debugClient *DebugClient) *Server {
	srv := &Server{
		router:             chi.NewRouter(),
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

	s.router.Get("/api/subs/last", s.fetchLastSub())
	s.router.Get("/subs/last", s.lastSubView(baseURL+"/api/subs/last"))

	s.router.Get("/api/gift/last", s.fetchLastGiftSub())
	s.router.Get("/gift/last", s.lastGiftSubView(baseURL+"/api/gift/last"))

	s.router.Get("/api/bits/last", s.fetchLastBits())
	s.router.Get("/bits/last", s.lastBitsView(baseURL+"/api/bits/last"))

	// DEBUG - trigger various events for testing
	// TODO how do secure when deploy?
	s.router.Get("/debug/sub", s.debugSub(s.debugClient))
	s.router.Get("/debug/gift", s.debugGift(s.debugClient))
	s.router.Get("/debug/bit", s.debugBit(s.debugClient))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
