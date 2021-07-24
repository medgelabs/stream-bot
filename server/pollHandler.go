package server

import (
	"encoding/json"
	"io"
	"medgebot/logger"
	"net/http"
	"time"
)

// currentPollView renders and returns the Poll on-screen HTML box
func (s *Server) currentPollView(apiEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := RefreshingView{
			ApiEndpoint: apiEndpoint,
		}
		s.pollHTML.Execute(w, data)
	}
}

// fetchCurrentPoll returns the current poll's question and voted answers state
func (s *Server) fetchCurrentPoll() http.HandlerFunc {
	type Answer struct {
		Label string `json:"label"`
		Count int    `json:"count"`
	}
	type response struct {
		Question string   `json:"question"`
		Answers  []Answer `json:"answers"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if !s.bot.IsPollRunning() {
			s.WriteJSON(w, 200, response{
				Question: "",
				Answers:  []Answer{},
			})

			return
		}

		question, answers := s.bot.GetPollState()
		resp := response{
			Question: question,
		}
		for _, answer := range answers {
			resp.Answers = append(resp.Answers, Answer{
				Label: answer.Answer,
				Count: answer.Count,
			})
		}

		s.WriteJSON(w, 200, resp)
	}
}

func (s *Server) createPoll() http.HandlerFunc {
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
		s.bot.StartPoll(time.Duration(req.Minutes)*time.Minute, req.Question, req.Answers)
	}
}
