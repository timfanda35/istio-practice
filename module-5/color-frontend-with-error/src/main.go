package main

import (
	  "fmt"
		"log"
		"net/http"
		"os"
)

const colorEndpoint = "http://color-backend/color"
const pageTitle = "Developing"
const indexTemple = `
<html>
  <head><title>{{.Title}}</title></head>
	<body>
	  <div style="text-align: center;">
			<h1 style="color: {{.Color}};">{{.Title}}</h1>

			<img src="https://placekitten.com/300/300" />
		</div>
	</body>
</html>
`
// DynamicContent is for page render
type DynamicContent struct {
	Title string
	Color string
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
  w.Write([]byte("500 - Something bad happened!"))
}

func main() {
    http.HandleFunc("/", index)
    port := os.Getenv("PORT")
		if port == "" {
						port = "8080"
						log.Printf("Defaulting to port %s", port)
		}

		log.Printf("Listening on port %s", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}