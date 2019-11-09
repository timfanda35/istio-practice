package main

import (
	  "fmt"
		"io/ioutil"
		"html/template"
		"log"
		"net/http"
		"os"
		"time"
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
	time.Sleep(10 * time.Second) // Delay 5 second

	resp, err := http.Get(colorEndpoint)
	if err != nil {
		log.Fatal("Call backend fail...", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Read color fail...", err)
	}

  content := DynamicContent{
		Title: pageTitle,
		Color: string(body),
	}
	tmpl, err := template.New("name").Parse(indexTemple)
	err = tmpl.Execute(w, content)
  if err != nil {
		log.Fatal("Render color fail...", err)
	}
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