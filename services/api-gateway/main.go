package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)

func main() {
	log.Println("Starting API Gateway")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /trip/preview", handleTripPreview)
	mux.HandleFunc("/ws/drivers", handleDriversWebSocket)
	mux.HandleFunc("/ws/riders", handleRidersWebSocket)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
	errChannel := make(chan error, 1)
	go func() {
		errChannel <- server.ListenAndServe()
		log.Printf("API Gateway is running on http://0.0.0.0:%v", httpAddr)
	}()
	shotDown := make(chan os.Signal, 1)
	signal.Notify(shotDown, os.Interrupt, syscall.SIGTERM)

	select {

	case err := <-errChannel:
		log.Fatalf("Could not start API Gateway: %v", err)
	case sig := <-shotDown:
		log.Printf("Shutting down API Gateway... Reason: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v", err)
		}
		log.Println("API Gateway stopped")
	}
}
