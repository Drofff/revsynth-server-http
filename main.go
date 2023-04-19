package main

import (
	"log"
	"net/http"

	"github.com/Drofff/revsynth-server-http/handler"
)

const serverAddr = ":8080"

func main() {
	log.Println("starting HTTP server at", serverAddr)
	http.HandleFunc("/api/v1/synth", handler.HandleRequest)
	err := http.ListenAndServe(serverAddr, http.DefaultServeMux)
	if err != nil {
		log.Fatalln("failed to start the server", err)
	}
}
