package server

import (
	"encoding/json"
	"html/template"
	"medgebot/bot/viewer"
	"medgebot/logger"
	log "medgebot/logger"
	"net/http"
	"sync"
)

// RefreshingView represents a view that polls for data to be interpolated
// on the View template
type RefreshingView struct {
	ApiEndpoint string
}

// fetchLastSub for the lastSubView
func (s *Server) fetchLastSub() http.HandlerFunc {
	type response struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		str, _ := s.viewerMetricsStore.Get(viewer.LastSub)
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
	var (
		onlyOnce sync.Once
		tmpl     *template.Template = template.New("lastSubsTemplate")
		err      error
	)

	return func(w http.ResponseWriter, r *http.Request) {
		onlyOnce.Do(func() {
			tmpl, err = tmpl.Parse(`
				<html lang="en">
				<head></head>
				<body>
					<section class="content"></section>

					 <script type="text/javascript">
					   fetchContent()
					   setInterval(fetchContent, 3000)

					   function fetchContent() {
						  let content = document.querySelector("section.content")
						  fetch("{{.ApiEndpoint}}")
						    .then(r => r.json())
							.then(r => content.innerHTML = "Last Sub: " + r.name)
							.catch(err => content.innerHTML = err)
						}
					 </script>
				</body>
				</html>`)
		})

		if err != nil {
			log.Fatal(err, "Last Subs template did not parse")
		}

		data := RefreshingView{
			ApiEndpoint: apiEndpoint,
		}
		tmpl.Execute(w, data)
	}
}

// fetchLastGiftSub for the lastGiftSubView
func (s *Server) fetchLastGiftSub() http.HandlerFunc {
	type response struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		str, _ := s.viewerMetricsStore.Get(viewer.LastGiftSub)
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
	var (
		onlyOnce sync.Once
		tmpl     *template.Template = template.New("lastGiftSubsTemplate")
		err      error
	)

	return func(w http.ResponseWriter, r *http.Request) {
		onlyOnce.Do(func() {
			tmpl, err = tmpl.Parse(`
				<html lang="en">
				<head>

				</head>
				<body>
					<section class="content"></section>

					 <script type="text/javascript">
					   fetchContent()
					   setInterval(fetchContent, 3000)

					   function fetchContent() {
						  let content = document.querySelector("section.content")
						  fetch("{{.ApiEndpoint}}")
						    .then(r => r.json())
							.then(r => content.innerHTML = "Last Gifter: " + r.name)
							.catch(err => content.innerHTML = err)
						}
					 </script>
				</body>
				</html>`)
		})

		if err != nil {
			log.Fatal(err, "Last Gift Subs template did not parse")
		}

		data := RefreshingView{
			ApiEndpoint: apiEndpoint,
		}
		tmpl.Execute(w, data)
	}
}
