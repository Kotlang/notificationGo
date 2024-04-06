package extensions

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type RequestBody struct {
	DeliveryInfoNotification struct {
		Subtid       string `json:"subtid"`
		DeliveryInfo struct {
			DeliveryChannel string `json:"deliveryChannel"`
			Description     string `json:"Description"`
			DestinationType string `json:"destinationType"`
			TimeStamp       string `json:"timeStamp"`
			Code            string `json:"code"`
			AdditionalInfo  string `json:"additionalInfo"`
			DeliveryStatus  string `json:"deliveryStatus"`
			Destination     string `json:"destination"`
			IdentityKeyHash string `json:"identityKeyHash"`
		} `json:"deliveryInfo"`
		CorrelationID string `json:"correlationid"`
		CallbackData  string `json:"callbackData"`
		TransID       string `json:"transid"`
	} `json:"deliveryInfoNotification"`
}

type TrimmedDeliveryInfo struct {
	Description    string
	DeliveryStatus string
	Recipient      string
	TimeStamp      int64
}

func (t TrimmedDeliveryInfo) GetType() string {
	return "delivery"
}

// DeliveryHandler handles requests for Deliver, Read and failed messages from the WhatsApp API
func DeliveryHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received delivery notification: /whatapp/delivery")
	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deliveryInfo := req.DeliveryInfoNotification

	info := TrimmedDeliveryInfo{
		Description:    deliveryInfo.DeliveryInfo.Description,
		DeliveryStatus: deliveryInfo.DeliveryInfo.DeliveryStatus,
		Recipient:      deliveryInfo.DeliveryInfo.Destination,
	}

	t, err := time.Parse(time.RFC3339, deliveryInfo.DeliveryInfo.TimeStamp)
	if err != nil {
		logger.Error("Error parsing timestamp", zap.Error(err))
	} else {
		info.TimeStamp = t.Unix()
	}

	trimmedTransactionId := strings.Split(deliveryInfo.TransID, "_")[0]

	WhatsappCache.Add(trimmedTransactionId, info)
}

type ReplyInfo struct {
	UserID          string    `json:"userId"`
	Username        string    `json:"username"`
	Channel         string    `json:"channel"`
	AppID           string    `json:"appId"`
	Event           string    `json:"event"`
	Waid            string    `json:"waid"`
	Timestamp       time.Time `json:"ts"`
	TransactionID   string    `json:"tid"`
	ButtonPayload   string    `json:"buttonPayload"`
	ButtonText      string    `json:"buttonText"`
	Errors          string    `json:"errors"`
	IdentityKeyHash string    `json:"identityKeyHash"`
}

func (w ReplyInfo) GetType() string {
	return "buttonReply"
}

func ReplyHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received reply notification: /whatsapp/reply")
	var req ReplyInfo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trimmedTransactionId := strings.Split(req.TransactionID, "_")[0]

	WhatsappCache.Add(trimmedTransactionId, req)
}
