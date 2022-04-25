package server

import (
	"encoding/json"
	"html/template"
	"medgebot/bot"
	"medgebot/cache"
	"medgebot/eventsub"
	"medgebot/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Server REST API
type Server struct {
	localBaseURL   string
	router         *chi.Mux
	bot            *bot.Bot
	store          cache.Cache
	debugClient    *DebugClient
	eventSubClient eventsub.Client
	labelHTML      *template.Template
	pollHTML       *template.Template
}

// New returns a Server instance to be run with http.ListenAndServe()
func New(localBaseURL string, bot *bot.Bot, dataStore cache.Cache, debugClient *DebugClient, labelHTMLStr string, pollHTMLStr string) *Server {
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

	// Parse Poll HTML template for reuse by the Poll View Handlers
	tmpl, err = template.New("PolTemplate").Parse(pollHTMLStr)
	if err != nil {
		logger.Fatal(err, "Failed to parse Poll HTML template")
	}
	srv.pollHTML = tmpl

	srv.routes()
	return srv
}

func (s *Server) routes() {
	// REQUIRED for EventSub callbacks
	s.router.Post("/eventsub/callback", s.eventSubHandler(s.eventSubClient))

	s.router.Get("/api/subs/last", s.fetchLastSub())
	s.router.Get("/subs/last", s.lastSubView(s.localBaseURL+"/api/subs/last"))

	s.router.Get("/api/gift/last", s.fetchLastGiftSub())
	s.router.Get("/gift/last", s.lastGiftSubView(s.localBaseURL+"/api/gift/last"))

	s.router.Get("/api/bits/last", s.fetchLastBits())
	s.router.Get("/bits/last", s.lastBitsView(s.localBaseURL+"/api/bits/last"))

	// Polls
	s.router.Post("/poll", s.createPoll())
	s.router.Get("/poll", s.currentPollView(s.localBaseURL+"/api/poll"))
	s.router.Get("/api/poll", s.fetchCurrentPoll())

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
