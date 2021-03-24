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

// fetchLastBits for the lastBitsView
func (s *Server) fetchLastBits() http.HandlerFunc {
	type response struct {
		Name   string `json:"name"`
		Months int    `json:"months"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		str, _ := s.viewerMetricsStore.Get(viewer.LastBits)
		lastBits, err := viewer.FromString(str)

		if err != nil {
			logger.Error(err, "lastBits cache fetch")
			w.WriteHeader(http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(response{
			Name:   lastBits.Name,
			Months: lastBits.Amount,
		})
	}
}

// lastBitsView returns the refreshing HTML page to grab the last Bitsscriber
func (s *Server) lastBitsView(apiEndpoint string) http.HandlerFunc {
	var (
		onlyOnce sync.Once
		tmpl     *template.Template = template.New("lastBitsTemplate")
		err      error
	)

	return func(w http.ResponseWriter, r *http.Request) {
		onlyOnce.Do(func() {
			// TODO can this template be part of a Base template somewhere?
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
							.then(r => content.innerHTML = "Last Bits: " + r.name)
							.catch(err => content.innerHTML = err)
						}
					 </script>
				</body>
				</html>`)
		})

		if err != nil {
			log.Fatal(err, "Last Bits template did not parse")
		}

		data := RefreshingView{
			ApiEndpoint: apiEndpoint,
		}
		tmpl.Execute(w, data)
	}
}
