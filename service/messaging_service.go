package service

import (
	"context"

	"github.com/Kotlang/notificationGo/clients"
	"github.com/Kotlang/notificationGo/db"
	notificationPb "github.com/Kotlang/notificationGo/generated/notification"
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/auth"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MessagingServiceInterface interface {
	BroadcastMessage(context.Context, *notificationPb.MesssageRequest) (*notificationPb.StatusResponse, error)
	RegisterMessagingTemplate(context.Context, *notificationPb.MessagingTemplate) (*notificationPb.StatusResponse, error)
	GetMessagingTemplate(context.Context, *notificationPb.IdRequest) (*notificationPb.MessagingTemplate, error)
}

type MessagingService struct {
	notificationPb.UnimplementedMessagingServiceServer
	db              db.NotificationDbInterface
	messagingClient clients.MessagingClient
}

func NewMessagingService(messagingClient clients.MessagingClient, db db.NotificationDbInterface) *MessagingService {
	return &MessagingService{
		db:              db,
		messagingClient: messagingClient,
	}
}

func (s *MessagingService) BroadcastMessage(ctx context.Context, req *notificationPb.MesssageRequest) (*notificationPb.StatusResponse, error) {
	_, tenant := auth.GetUserIdAndTenant(ctx)

	// Get the template
	exists := s.db.MessagingTemplate(tenant).IsExistsById(req.TemplateId)
	if !exists {
		return nil, status.Error(codes.NotFound, "Template not found")
	}

	parameters := getParameter(req.MediaParameters, req.TemplateParameters)

	// Send message to the destination
	_, err := s.messagingClient.SendMessage(req.TemplateId, req.RecipientPhoneNumber, parameters)
	if err != nil {
		return nil, err
	}

	return &notificationPb.StatusResponse{
		Status: "success",
	}, nil
}

func (s *MessagingService) RegisterMessagingTemplate(ctx context.Context, req *notificationPb.MessagingTemplate) (*notificationPb.StatusResponse, error) {
	_, tenant := auth.GetUserIdAndTenant(ctx)

	// Register the template
	messagingTemplateModel := getMessagingTemplateModel(req)
	err := <-s.db.MessagingTemplate(tenant).Save(messagingTemplateModel)
	if err != nil {
		return nil, err
	}

	return &notificationPb.StatusResponse{
		Status: "success",
	}, nil
}

func getMessagingTemplateModel(req *notificationPb.MessagingTemplate) *models.MessagingTemplateModel {
	messagingTemplateModel := &models.MessagingTemplateModel{}
	copier.CopyWithOption(messagingTemplateModel, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	return messagingTemplateModel
}

func getParameter(mediaParameters *notificationPb.MediaParameters, templateParameters map[string]string) map[string]interface{} {
	parameters := make(map[string]interface{})
	for k, v := range templateParameters {
		parameters[k] = v
	}

	if mediaParameters == nil {
		return parameters
	}

	mediaType := mediaParameters.MediaType.String()
	parameters[mediaType] = map[string]string{
		"link":     mediaParameters.Link,
		"filename": mediaParameters.Filename,
	}
	return parameters
}
