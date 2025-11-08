package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
)

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

	// Marshal the parsed request body and send that to the trip service
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}

	TripResponse, err := http.Post("http://trip-service:8083/preview", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "Failed to communicate with trip service", http.StatusInternalServerError)
		return
	}
	defer TripResponse.Body.Close()

	var response contracts.APIResponse
	if err := json.NewDecoder(TripResponse.Body).Decode(&response); err != nil {
		http.Error(w, "Invalid response from trip service", http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, http.StatusOK, response)

}
