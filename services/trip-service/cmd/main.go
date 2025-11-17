package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"
	"time"
)

func main() {
	inmemRepo := repository.NewInmemRepository()
	svc := service.NewService(inmemRepo)
	mux := http.NewServeMux()

	httphandler := h.HttpHandler{Service: svc}

	mux.HandleFunc("POST /preview", httphandler.HandleTripPreview)

	server := &http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	errChannel := make(chan error, 1)
	go func() {
		errChannel <- server.ListenAndServe()
		log.Printf("Trip Service is running on http://0.0.0.0.0:", server.Addr)
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-errChannel:
		log.Fatalf("could not start server: %v", err)

	case sig := <-shutdown:
		log.Printf("shutting down server... reason: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("could not gracefully shutdown the server: %v", err)
		}
		log.Printf("server stopped Successfully")

	}
}
