package main

import (
	"fmt"
	"log"
	"net/http"

	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {

	log.Println("Starting Trip Service...")
	mux := http.NewServeMux()

	mux.HandleFunc("POST /preview", handleTripPreview)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
	log.Printf("Trip Service is listening on %s", httpAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
	fmt.Printf("server started on %v", httpAddr)

}
