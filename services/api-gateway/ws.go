package main

import (
	"fmt"
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"
)

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
		return
	}
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")

	if packageSlug == "" {
		http.Error(w, "packageSlug is required", http.StatusBadRequest)
		return
	}
	type Driver struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ProfilePic  string `json:"profilePicture"`
		CarPlate    string `json:"carPlate"`
		PackageSlug string `json:"packageSlug"`
	}
	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID:          userID,
			Name:        "John Doe",
			ProfilePic:  util.GetRandomAvatar(1234),
			CarPlate:    "XYZ 1234",
			PackageSlug: packageSlug,
		},
	}
	defer conn.Close()

	err = conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		log.Printf("Received message: %s", msg)
	}

}
