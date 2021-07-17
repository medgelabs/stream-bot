package server

import (
	"encoding/json"
	"io"
	"medgebot/logger"
	"net/http"
	"time"
)

func (s *Server) createPoll() func(http.ResponseWriter, *http.Request) {
	type request struct {
		Question string   `json:"question"`
		Answers  []string `json:"answers"`
		Minutes  int      `json:"minutes,omitempty"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			body, _ := io.ReadAll(r.Body)
			defer r.Body.Close()
			logger.Error(err, "Failed to unmarshal creatPoll request: %s", body)
			s.WriteError(w, 400, "Invalid request body")
			return
		}
		defer r.Body.Close()

		if req.Question == "" {
			s.WriteError(w, 400, "request body missing Question")
			return
		}

		if len(req.Answers) == 0 {
			s.WriteError(w, 400, "request body needs at least 1 answer")
			return
		}

		// TODO accept on the request
		req.Minutes = 3

		// Check if poll already running. If yes - return error
		if s.bot.IsPollRunning() {
			s.WriteError(w, 409, "Poll already running")
			return
		}

		// > write poll to cache, trigger Bot into Poll mode
		s.store.Clear("pollAnswers")
		s.bot.StartPoll(time.Duration(req.Minutes)*time.Minute, req.Question, req.Answers)
	}
}
