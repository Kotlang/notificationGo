package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

const apiURL = "https://api.imiconnect.io/resources/v1/messaging"

type WhatsAppClient struct {
	appId, serviceKey string
	client            *http.Client
}

type ResponseBody struct {
	Response []ResponseData `json:"response"`
}

type ResponseData struct {
	Code          string `json:"code"`
	TransactionID string `json:"transid"`
	Description   string `json:"description"`
	CorrelationID string `json:"correlationid"`
}

var whatsappClient *WhatsAppClient

func NewWhatsAppClient() *WhatsAppClient {

	if whatsappClient == nil {
		appId := os.Getenv("IMI_APP_ID")
		serviceKey := os.Getenv("IMI_SERVICE_KEY")

		if appId == "" || serviceKey == "" {
			logger.Error("IMI_APP_ID or IMI_SERVICE_KEY is not set")
			return nil
		}

		httpClient := &http.Client{}
		whatsappClient = &WhatsAppClient{appId, serviceKey, httpClient}
	}

	return whatsappClient
}

// SendMessage sends a message to the given destination and returns the transaction ID.
func (w *WhatsAppClient) SendMessage(templateID string, destination []string, parameters map[string]interface{}) (transactionId string, err error) {
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

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body", zap.Error(err))
		return "", err
	}

	// unmarshal the response body
	var respBody ResponseBody
	err = json.Unmarshal(responseBody, &respBody)
	if err != nil {
		logger.Error("Error unmarshalling response body", zap.Error(err))
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		logger.Error("Failed to send message", zap.String("status", response.Status))
		return "", fmt.Errorf("failed to send message: %s", response.Status)
	}

	return respBody.Response[0].TransactionID, nil
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
