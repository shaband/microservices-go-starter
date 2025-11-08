package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/contracts"
)

var (
	ctx = context.Background()

	inmemRepo = repository.NewInmemRepository()

	svc = service.NewService(inmemRepo)
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Trip Preview Request")
	route, err := svc.CreateRoute(ctx)
	if err != nil {
		http.Error(w, "Failed to Create Route", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := contracts.APIResponse{
		Data: route,
	}

	json.NewEncoder(w).Encode(response)
}
