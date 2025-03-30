package main

import (
	"log"
	"net/http"
	"ws-chat/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWs()

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("Error run web server: %s", err)
		return
	}
}
