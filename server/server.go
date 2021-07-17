package server

import (
	"encoding/json"
	"html/template"
	"medgebot/bot"
	"medgebot/cache"
	"medgebot/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server REST API
type Server struct {
	router      *chi.Mux
	bot         *bot.Bot
	store       cache.Cache
	debugClient *DebugClient
	labelHTML   *template.Template
}

// New returns a Server instance to be run with http.ListenAndServe()
func New(bot *bot.Bot, dataStore cache.Cache, debugClient *DebugClient, labelHTMLStr string) *Server {
	srv := &Server{
		router:      chi.NewRouter(),
		store:       dataStore,
		debugClient: debugClient,
		bot:         bot,
	}

	// Parse Metrics Label HTML template for reuse by the View Handlers
	tmpl, err := template.New("MetricsLabel").Parse(labelHTMLStr)
	if err != nil {
		logger.Fatal(err, "Failed to parse Metrics Label HTML template")
	}
	srv.labelHTML = tmpl

	// TODO receive WebSocket connection for alerts and assign to alertWebSocket

	srv.routes()
	return srv
}

func (s *Server) routes() {
	// TODO pull from config
	baseURL := "http://localhost:8080"

	// Metrics endpoints
	s.router.Get("/api/subs/last", s.fetchLastSub())
	s.router.Get("/subs/last", s.lastSubView(baseURL+"/api/subs/last"))

	s.router.Get("/api/gift/last", s.fetchLastGiftSub())
	s.router.Get("/gift/last", s.lastGiftSubView(baseURL+"/api/gift/last"))

	s.router.Get("/api/bits/last", s.fetchLastBits())
	s.router.Get("/bits/last", s.lastBitsView(baseURL+"/api/bits/last"))

	// Polls
	s.router.Post("/polls", s.createPoll())

	// DEBUG - trigger various events for testing
	// TODO how do secure when deploy?
	s.router.Get("/debug/sub", s.debugSub(s.debugClient))
	s.router.Get("/debug/gift", s.debugGift(s.debugClient))
	s.router.Get("/debug/bit", s.debugBit(s.debugClient))
}

// WriteJSON is a helper to respond with a JSON message body.
// If marshalling fails, it will respond with a HTTP 500
func (s *Server) WriteJSON(w http.ResponseWriter, statusCode int, msg interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		logger.Error(err, "Failed to write Error response")
		s.WriteError(w, 500, "Response marshalling failed")
	}
}

// Standard errors
type error struct {
	Error string `json:"error"`
}

// WriteError responds with a standardized error body
func (s *Server) WriteError(w http.ResponseWriter, statusCode int, msg string) {
	s.WriteJSON(w, statusCode, error{
		Error: msg,
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
