package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type MessageRepositoryInterface interface {
	odm.BootRepository[models.MessageModel]
	GetMessagesByReceiver(receiver string, limit, skip int64) []models.MessageModel
	GetMessagesBySender(sender string, limit, skip int64) []models.MessageModel
}

type MessageRepository struct {
	odm.UnimplementedBootRepository[models.MessageModel]
}

func (m *MessageRepository) GetMessagesByReceiver(receiver string, limit, skip int64) []models.MessageModel {

	filters := bson.M{"receiver": receiver}
	sort := bson.D{{Key: "createdOn", Value: -1}}

	resultChan, errChan := m.Find(filters, sort, limit, skip)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting messages", zap.Error(err))
		return nil
	}
}

func (m *MessageRepository) GetMessagesBySender(sender string, limit, skip int64) []models.MessageModel {

	filters := bson.M{"sender": sender}

	sort := bson.D{{Key: "createdOn", Value: -1}}

	resultChan, errChan := m.Find(filters, sort, limit, skip)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting messages", zap.Error(err))
		return nil
	}
}

func (m *MessageRepository) GetUserChatHistory(user string, limit, skip int64) []models.MessageModel {

	filters := bson.M{"$or": []bson.M{{"sender": user}, {"receiver": user}}}

	sort := bson.D{{Key: "createdOn", Value: -1}}

	resultChan, errChan := m.Find(filters, sort, limit, skip)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting messages", zap.Error(err))
		return nil
	}
}
