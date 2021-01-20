package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

// RefreshingView represents a view that polls for data to be interpolated
// on the View template
type RefreshingView struct {
	ApiEndpoint string
}

// fetchLastSub for the lastSubView
func (s *server) fetchLastSub() http.HandlerFunc {
	type response struct {
		Name   string `json:"name"`
		Months int    `json:"months"`
	}

	// TODO hook into Metrics Store
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(response{
			Name:   fmt.Sprintf("%d", time.Now().Unix()),
			Months: 3,
		})
	}
}

// lastSubView returns the refreshing HTML page to grab the last Subscriber
func (s *server) lastSubView(apiEndpoint string) http.HandlerFunc {
	var (
		onlyOnce sync.Once
		tmpl     *template.Template = template.New("lastSubsTemplate")
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
							.then(r => content.innerHTML = r.name)
							.catch(err => content.innerHTML = err)
						}
					 </script>
				</body>
				</html>`)
		})

		if err != nil {
			log.Fatalf("FATAL: Last Subs template did not parse - %v", err)
		}

		data := RefreshingView{
			ApiEndpoint: apiEndpoint,
		}
		tmpl.Execute(w, data)
	}
}
