package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type MessagingTemplateRepositoryInterface interface {
	odm.BootRepository[models.MessagingTemplateModel]
	GetTemplateByType(templateType string, limit, skip int64) []models.MessagingTemplateModel
	GetTemplateByName(templateName string) *models.MessagingTemplateModel
}

type MessagingTemplateRepository struct {
	odm.UnimplementedBootRepository[models.MessagingTemplateModel]
}

func (m *MessagingTemplateRepository) GetTemplateByType(templateType string, limit, skip int64) []models.MessagingTemplateModel {

	filters := bson.M{"templateType": templateType}

	resultChan, errChan := m.Find(filters, bson.D{}, limit, skip)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting templates", zap.Error(err))
		return nil
	}
}

func (m *MessagingTemplateRepository) GetTemplateByName(templateName string) *models.MessagingTemplateModel {

	filters := bson.M{"templateName": templateName}

	resultChan, errChan := m.FindOne(filters)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting template", zap.Error(err))
		return nil
	}
}
