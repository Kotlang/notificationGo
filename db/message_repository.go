package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	notificationPb "github.com/Kotlang/notificationGo/generated/notification"
)

type MessageRepositoryInterface interface {
	odm.BootRepository[models.MessageModel]
	GetMessageByTransactionId(transactionId string) *models.MessageModel
	GetMessages(messageFilters *notificationPb.MessageFilters, pageNumber, pageSize int64) (messages []models.MessageModel, totalCount int64)
}

type MessageRepository struct {
	odm.UnimplementedBootRepository[models.MessageModel]
}

func (m *MessageRepository) GetUserChatHistory(user string, limit, skip int64) []models.MessageModel {

	filter := bson.M{"$or": []bson.M{{"sender": user}, {"receiver": user}}}

	sort := bson.D{{Key: "createdOn", Value: -1}}

	resultChan, errChan := m.Find(filter, sort, limit, skip)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting messages", zap.Error(err))
		return nil
	}
}

func (m *MessageRepository) GetMessageByTransactionId(transactionId string) *models.MessageModel {

	filters := bson.M{"transactionId": transactionId}

	resultChan, errChan := m.FindOne(filters)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting message", zap.Error(err))
		return nil
	}
}

func (m *MessageRepository) GetMessages(messageFilters *notificationPb.MessageFilters, pageNumber, pageSize int64) (messages []models.MessageModel, totalCount int64) {

	if pageNumber < 0 {
		pageNumber = 0
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	filter := bson.M{}

	if messageFilters != nil {
		if messageFilters.Sender != "" {
			filter["sender"] = messageFilters.Sender
		}

		if messageFilters.TransactionId != "" {
			filter["transactionId"] = messageFilters.TransactionId
		}
	}

	skip := pageNumber * pageSize

	sort := bson.D{{Key: "createdOn", Value: -1}}

	// fetch messages from db
	totalCountResChan, countErrChan := m.CountDocuments(filter)
	messagesResChan, messageErrChan := m.Find(filter, sort, pageSize, skip)

	select {
	case totalCount = <-totalCountResChan:
	case err := <-countErrChan:
		logger.Error("Failed getting messages count", zap.Error(err))
		return nil, 0
	}

	select {
	case messages = <-messagesResChan:
		return messages, totalCount
	case err := <-messageErrChan:
		logger.Error("Failed getting messages", zap.Error(err))
		return nil, 0
	}

}
