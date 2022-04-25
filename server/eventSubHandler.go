package server

import (
	"medgebot/eventsub"
	"medgebot/logger"
	"net/http"
)

func (s *Server) eventSubHandler(client eventsub.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO switch for challenge key
		challengeKey := client.Secret()
		logger.Info("Responding to Event Sub challenge with %s", challengeKey)

		w.WriteHeader(200)
		w.Write([]byte(challengeKey))

		// TODO parse events
	}
}
