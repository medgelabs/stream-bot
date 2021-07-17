package server

import (
	"encoding/json"
	"medgebot/bot/viewer"
	"medgebot/logger"
	"net/http"
)

// RefreshingView represents a view that polls for data to be interpolated
// on the View template
type RefreshingView struct {
	ApiEndpoint string
	Label       string
}

// fetchLastSub for the lastSubView
// Note: response must be a `data` field for the common metric HTML template
func (s *Server) fetchLastSub() http.HandlerFunc {
	type response struct {
		Name string `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		str, _ := s.store.Get(viewer.LastSub)
		lastSub, err := viewer.FromString(str)

		if err != nil {
			logger.Error(err, "lastSub cache fetch")
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(response{
			Name: lastSub.Name,
		})
	}
}

// lastSubView returns the refreshing HTML page to grab the last Subscriber
func (s *Server) lastSubView(apiEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := RefreshingView{
			ApiEndpoint: apiEndpoint,
			Label:       "Last Sub",
		}
		s.labelHTML.Execute(w, data)
	}
}

// fetchLastGiftSub for the lastGiftSubView
// Note: response must be a `data` field for the common metric HTML template
func (s *Server) fetchLastGiftSub() http.HandlerFunc {
	type response struct {
		Name string `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		str, _ := s.store.Get(viewer.LastGiftSub)
		lastGifter, err := viewer.FromString(str)

		if err != nil {
			logger.Error(err, "lastGiftSub cache fetch")
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(response{
			Name: lastGifter.Name,
		})
	}
}

// lastGiftSubView returns the refreshing HTML page to grab the last Gifter
func (s *Server) lastGiftSubView(apiEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := RefreshingView{
			ApiEndpoint: apiEndpoint,
			Label:       "Last Gifter",
		}
		s.labelHTML.Execute(w, data)
	}
}
