package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	notificationPb "github.com/Kotlang/notificationGo/generated/notification"
)

type MessagingTemplateRepositoryInterface interface {
	odm.BootRepository[models.MessagingTemplateModel]
	GetTemplate(notificationFilters *notificationPb.FetchTemplateRequest) (templates []models.MessagingTemplateModel, totalCount int64)
}

type MessagingTemplateRepository struct {
	odm.UnimplementedBootRepository[models.MessagingTemplateModel]
}

func (m *MessagingTemplateRepository) GetTemplate(notificationFilter *notificationPb.FetchTemplateRequest) (templates []models.MessagingTemplateModel, totalCount int64) {

	filter := bson.M{}

	if notificationFilter.TemplateId != "" {
		filter["_id"] = notificationFilter.TemplateId
	}

	if notificationFilter.TemplateName != "" {
		filter["templateName"] = notificationFilter.TemplateName
	}

	skip := notificationFilter.PageNumber * notificationFilter.PageSize

	sort := bson.D{{Key: "createdOn", Value: -1}}

	// fetch templates from db
	totalCountResChan, countErrChan := m.CountDocuments(filter)
	templatesResChan, templateErrChan := m.Find(filter, sort, int64(notificationFilter.PageSize), int64(skip))

	select {
	case totalCount = <-totalCountResChan:
	case err := <-countErrChan:
		logger.Error("Failed getting templates count", zap.Error(err))
		return nil, 0
	}

	select {
	case templates = <-templatesResChan:
		return templates, totalCount
	case err := <-templateErrChan:
		logger.Error("Failed getting templates", zap.Error(err))
		return nil, 0
	}

}
