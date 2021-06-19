package server

import (
	"medgebot/bot"
	"medgebot/cache"
	"medgebot/eventsub"
	"net/http"
)

// Server REST API
type Server struct {
	router             *http.ServeMux
	viewerMetricsStore cache.Cache
	eventSubClient     eventsub.Client
	alertWebSocket     *bot.WriteOnlyUnsafeWebSocket
	debugClient        *DebugClient
	localBaseURL       string
}

// New returns a Server instance to be run with http.ListenAndServe()
func New(localBaseURL string, metricStore cache.Cache, eventSubClient eventsub.Client, alertWebSocket *bot.WriteOnlyUnsafeWebSocket, debugClient *DebugClient) *Server {
	srv := &Server{
		router:             http.NewServeMux(),
		viewerMetricsStore: metricStore,
		eventSubClient:     eventSubClient,
		alertWebSocket:     alertWebSocket,
		debugClient:        debugClient,
		localBaseURL:       localBaseURL,
	}

	// TODO receive WebSocket connection for alerts and assign to alertWebSocket

	srv.routes()
	return srv
}

func (s *Server) routes() {
	// REQUIRED for EventSub callbacks
	s.router.HandleFunc("/eventsub/callback", s.eventSubHandler(s.eventSubClient))

	s.router.HandleFunc("/api/subs/last", s.fetchLastSub())
	s.router.HandleFunc("/subs/last", s.lastSubView(s.localBaseURL+"/api/subs/last"))

	s.router.HandleFunc("/api/gift/last", s.fetchLastGiftSub())
	s.router.HandleFunc("/gift/last", s.lastGiftSubView(s.localBaseURL+"/api/gift/last"))

	s.router.HandleFunc("/api/bits/last", s.fetchLastBits())
	s.router.HandleFunc("/bits/last", s.lastBitsView(s.localBaseURL+"/api/bits/last"))

	// DEBUG - trigger various events for testing
	// TODO how do secure when deploy?
	s.router.HandleFunc("/debug/sub", s.debugSub(s.debugClient))
	s.router.HandleFunc("/debug/gift", s.debugGift(s.debugClient))
	s.router.HandleFunc("/debug/bit", s.debugBit(s.debugClient))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
