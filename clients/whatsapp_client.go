package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

const apiURL = "https://api.imiconnect.io/resources/v1/messaging"

type WhatsAppClient struct {
	appId, serviceKey string
	client            *http.Client
}

func NewWhatsAppClient(appId, serviceKey string, client *http.Client) *WhatsAppClient {
	return &WhatsAppClient{appId, serviceKey, client}
}

func (w *WhatsAppClient) SendMessage(templateID string, destination []string, parameters map[string]interface{}) (string, error) {
	payload := getPayload(w.appId, templateID, destination, parameters)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Error marshalling payload", zap.Error(err))
		return "", err
	}

	request, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error("Error creating request", zap.Error(err))
		return "", err
	}

	// Set headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("key", w.serviceKey)

	response, err := w.client.Do(request)
	if err != nil {
		logger.Error("Error sending request", zap.Error(err))
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to send message: %s", response.Status)
	}

	return "Message sent successfully", nil
}

func getPayload(appID, templateID string, destination []string, parameters map[string]interface{}) map[string]interface{} {
	payload := map[string]interface{}{
		"appid":           appID,
		"deliverychannel": "whatsapp",
		"message": map[string]interface{}{
			"template":   templateID,
			"parameters": parameters,
		},
		"destination": []map[string]interface{}{
			{"waid": destination},
		},
	}
	return payload
}
