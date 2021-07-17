package server

import (
	"encoding/json"
	"medgebot/bot/viewer"
	"medgebot/logger"
	"net/http"
)

// fetchLastBits for the lastBitsView
func (s *Server) fetchLastBits() http.HandlerFunc {
	type response struct {
		Name string `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		str, _ := s.store.Get(viewer.LastBits)
		lastBits, err := viewer.FromString(str)

		if err != nil {
			logger.Error(err, "lastBits cache fetch")
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(response{
			Name: lastBits.Name,
		})
	}
}

// lastBitsView returns the refreshing HTML page to grab the last Bits
func (s *Server) lastBitsView(apiEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := RefreshingView{
			ApiEndpoint: apiEndpoint,
			Label:       "Last Bits",
		}
		s.labelHTML.Execute(w, data)
	}
}
