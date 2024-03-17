package service

import (
	"context"
	"fmt"

	"github.com/Kotlang/notificationGo/clients"
	"github.com/Kotlang/notificationGo/db"
	notificationPb "github.com/Kotlang/notificationGo/generated/notification"
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/auth"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
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
	templateResChan, errChan := s.db.MessagingTemplate(tenant).FindOneById(req.TemplateId)
	var template *models.MessagingTemplateModel
	select {
	case template = <-templateResChan:
	case err := <-errChan:
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("Template with id %s not found", req.TemplateId))
		}
		return nil, err
	}

	// TODO: Using Template generate message and save in the database
	logger.Info("Template: ", zap.Any("template", template))

	parameters := getParameter(req.MediaParameters, req.HeaderParameters, req.BodyParameters, req.ButtonParameters)
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
	model := &models.MessagingTemplateModel{}
	// Copy basic fields with copier
	copier.CopyWithOption(model, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// copy category
	model.Category = req.Category.String()

	// copy media parameters
	if req.MediaParameters != nil {
		model.MediaParameters.MediaType = req.MediaParameters.MediaType.String()
	}

	// copy button type
	model.ButtonType = req.ButtonType.String()

	// copy buttons
	if req.Buttons != nil && req.Buttons.CallToActionButtons != nil {
		buttons := make([]models.CallToActionButtons, 0)
		for _, button := range req.Buttons.CallToActionButtons {
			buttons = append(buttons, models.CallToActionButtons{
				ActionType:  button.ActionType.String(),
				Text:        button.Text,
				PhoneNumber: button.PhoneNumber,
				Url: models.Url{
					UrlType:       button.Url.UrlType.String(),
					Link:          button.Url.Link,
					UrlParameters: button.Url.UrlParameters,
				},
			})
		}
		model.Buttons.CallToActionButtons = buttons
	}

	return model
}

func getParameter(mediaParameters *notificationPb.MediaParameters, headerParameters, bodyParameters, buttonParameters map[string]string) map[string]interface{} {
	parameters := make(map[string]interface{})
	for k, v := range headerParameters {
		parameters[k] = v
	}

	for k, v := range bodyParameters {
		parameters[k] = v
	}

	for k, v := range buttonParameters {
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
