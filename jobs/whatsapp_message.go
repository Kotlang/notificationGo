package jobs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Kotlang/notificationGo/clients"
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/models"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type whatsappMessage struct {
	Name            string
	db              db.NotificationDbInterface
	messagingClient clients.MessagingClientInterface
}

func NewWhatsappMessageJob(db db.NotificationDbInterface) *whatsappMessage {
	return &whatsappMessage{
		Name:            "whatsapp.message",
		db:              db,
		messagingClient: clients.NewWhatsAppClient(),
	}
}

func (j *whatsappMessage) Run() (err error) {

	filter := bson.M{
		"eventType": j.Name,
		"templateParameters.scheduleTime": bson.M{
			"$lte": strconv.Itoa(int(time.Now().Unix())),
		},
	}

	eventsChan, errChan := j.db.Event().Find(filter, nil, 10, 0)

	var events []models.EventModel
	select {
	case events = <-eventsChan:
		fmt.Println("events", events)
	case err = <-errChan:
		logger.Error("Failed getting events", zap.Error(err))
		return
	}

	if len(events) == 0 {
		return
	}

	for _, event := range events {

		// if no target users delete the event log it and continue
		if len(event.TargetUsers) == 0 {
			logger.Error("Failed sending message to user",
				zap.Error(status.Error(codes.InvalidArgument, "no target users for event")))

			<-j.db.Event().DeleteById(event.Id())
			continue
		}

		// unmarshall the parameters field of templateParameters
		var parameters map[string]interface{}
		if err = json.Unmarshal([]byte(event.TemplateParameters["parameters"]), &parameters); err != nil {
			logger.Error("Failed unmarshalling template parameters", zap.Error(err))
			<-j.db.Event().DeleteById(event.Id())
			continue
		}

		transactionID, err := j.messagingClient.SendMessage(event.TemplateParameters["templateId"], event.TargetUsers, parameters)
		if err != nil {
			logger.Error("Error sending Message", zap.Error(err))
		}
		logger.Info("Succesfully sent message", zap.Any("response: %v", transactionID))

		// update transactionId of the message
		messageResChan, errChan := j.db.Message().FindOneById(event.TemplateParameters["messageId"])
		select {
		case message := <-messageResChan:
			message.TransactionId = transactionID
			err = <-j.db.Message().Save(message)
			if err != nil {
				logger.Error("Failed saving message", zap.Error(err))
			}
		case err = <-errChan:
			logger.Error("Failed getting message", zap.Error(err))
		}

		err = <-j.db.Event().DeleteById(event.Id())
		if err != nil {
			logger.Error("Failed deleting event", zap.Error(err))
		}
	}
	return err
}
