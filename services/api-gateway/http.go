package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/types"

	pb "ride-sharing/shared/proto/trip"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {

	var reqBody pb.PreviewTripRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if reqBody.UserID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	// Marshal the parsed request body and send that to the trip service
	bodyBytes, err := json.Marshal(&reqBody)
	if err != nil {
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}

	tripServiceClient, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		http.Error(w, "Failed to create trip service client", http.StatusInternalServerError)
		return
	}

	defer tripServiceClient.Close()
	tripServiceClient.Client.PreviewTrip(r.Context(), &reqBody)

	TripResponse, err := http.Post("http://trip-service:8083/preview", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "Failed to communicate with trip service", http.StatusInternalServerError)
		return
	}
	defer TripResponse.Body.Close()
	var response types.OsrmRespBody
	err = json.NewDecoder(TripResponse.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Invalid response from trip service", http.StatusInternalServerError)
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
