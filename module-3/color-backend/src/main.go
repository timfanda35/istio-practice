package main

import (
		"fmt"
		"log"
		"net/http"
		"os"
)

func color(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "green")
}

func main() {
    http.HandleFunc("/", color)
    port := os.Getenv("PORT")
		if port == "" {
						port = "8080"
						log.Printf("Defaulting to port %s", port)
		}

		log.Printf("Listening on port %s", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}