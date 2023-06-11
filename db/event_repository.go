package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type EventRepository struct {
	odm.AbstractRepository[models.EventModel]
}

func (e *EventRepository) GetEventByEventType(eventType string, limit, skip int64) []models.EventModel {

	filters := bson.M{"eventType": eventType}

	resultChan, errChan := e.Find(filters, bson.D{}, limit, skip)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting events", zap.Error(err))
		return []models.EventModel{}
	}
}
