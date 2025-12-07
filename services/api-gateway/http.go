package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"

	// pb "ride-sharing/shared/proto/trip"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {

	var reqBody previewTripRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if reqBody.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}
	tripServiceClient, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		http.Error(w, "Failed to create trip service client", http.StatusInternalServerError)
		return
	}

	defer tripServiceClient.Close()
	response, err := tripServiceClient.Client.PreviewTrip(r.Context(), reqBody.ToProto())
	if err != nil {
		log.Printf("Error previewing trip: %v", err)
		http.Error(w, "Failed to preview trip", http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, http.StatusOK, response)

}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)

	}

	defer conn.Close()
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\\n", message)

	}

}
