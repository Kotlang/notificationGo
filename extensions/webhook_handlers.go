package extensions

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IncominMessage struct {
	UserID          string `json:"userId"`
	Channel         string `json:"channel"`
	AppID           string `json:"appId"`
	Event           string `json:"event"`
	WaID            string `json:"waid"`
	Message         string `json:"message"`
	Attachments     string `json:"attachments"`
	System          string `json:"system"`
	Location        string `json:"location"`
	TS              string `json:"ts"`
	TID             string `json:"tid"`
	IdentityKeyHash string `json:"identityKeyHash"`
	Contacts        string `json:"contacts"`
}

func WhatsappIncomingMessageHandler(w http.ResponseWriter, r *http.Request) {
	var message IncominMessage

	// Parse the incoming JSON payload
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Handle different types of messages
	switch {
	case message.Message != "":
		fmt.Println("Received Text Message:", message.Message)
	case message.Attachments != "":
		// Check the type of attachment
		// For simplicity, let's assume the attachments field contains JSON data
		var attachments []map[string]string
		if err := json.Unmarshal([]byte(message.Attachments), &attachments); err != nil {
			fmt.Println("Failed to parse attachments:", err)
			return
		}
		for _, attachment := range attachments {
			mime, ok := attachment["mime_type"]
			if !ok {
				continue
			}
			switch mime {
			case "image/jpeg", "image/png":
				fmt.Println("Received Image:", attachment["url"])
			case "video/mp4":
				fmt.Println("Received Video:", attachment["url"])
			case "audio/mpeg":
				fmt.Println("Received Audio:", attachment["url"])
			case "application/msword":
				fmt.Println("Received Document:", attachment["url"])
			default:
				fmt.Println("Received Unknown Attachment Type:", mime)
			}
		}
	case message.Location != "":
		fmt.Println("Received Location:", message.Location)
	case message.Contacts != "":
		fmt.Println("Received Contact:", message.Contacts)
	default:
		fmt.Println("Received Unknown Message Type")
	}

	// Send a response if necessary
	w.WriteHeader(http.StatusOK)
}
