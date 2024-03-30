package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/Kotlang/notificationGo/clients"
	"github.com/Kotlang/notificationGo/db"
	notificationPb "github.com/Kotlang/notificationGo/generated/notification"
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/auth"
	"github.com/SaiNageswarS/go-api-boot/cloud"
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
	FetchMessagingTemplates(context.Context, *notificationPb.FetchTemplateRequest) (*notificationPb.MessagingTemplateList, error)
}

type MessagingService struct {
	notificationPb.UnimplementedMessagingServiceServer
	db              db.NotificationDbInterface
	messagingClient clients.MessagingClientInterface
	authClient      clients.AuthClientInterface
	cloudFns        cloud.Cloud
}

func NewMessagingService(messagingClient clients.MessagingClientInterface,
	db db.NotificationDbInterface,
	authClient clients.AuthClientInterface,
	cloudFns cloud.Cloud) *MessagingService {

	return &MessagingService{
		db:              db,
		messagingClient: messagingClient,
		authClient:      authClient,
		cloudFns:        cloudFns,
	}
}

func (s *MessagingService) BroadcastMessage(ctx context.Context, req *notificationPb.MesssageRequest) (*notificationPb.StatusResponse, error) {
	userId, tenant := auth.GetUserIdAndTenant(ctx)

	// check if user is admin
	if !<-s.authClient.IsUserAdmin(ctx, userId) {
		logger.Error("User is not admin", zap.String("userId", userId))
		return nil, status.Error(codes.PermissionDenied, "User is not admin")
	}
	//Get the template
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

	// TODO: Add validation for template and parameters recieved
	logger.Info("Template: ", zap.Any("template", template))

	// message model
	message := getMessageModel(userId, req)

	// get single map of all parameters
	parameters := getParameter(req.MediaParameters, req.HeaderParameters, req.BodyParameters, req.ButtonParameters)

	if req.ScheduleInfo != nil && req.ScheduleInfo.IsScheduled {

		// convert parameters to string
		parametersString, err := json.Marshal(parameters)
		if err != nil {
			logger.Error("Failed to marshal parameters", zap.Error(err))
			return nil, err
		}

		templateParams := make(map[string]string)

		templateParams["parameters"] = string(parametersString)
		templateParams["templateId"] = req.TemplateId
		templateParams["scheduleTime"] = fmt.Sprint(req.GetScheduleInfo().ScheduledTime)

		event := &models.EventModel{
			CreatorId:          userId,
			EventType:          "whatsapp.message",
			Title:              template.Header,
			Body:               template.Body,
			TemplateParameters: templateParams,
			Topic:              "message",
			TargetUsers:        req.RecipientPhoneNumber,
			Tenant:             tenant,
		}

		err = <-s.db.Event().Save(event)

		if err != nil {
			logger.Error("Failed to schedule message event", zap.Error(err))
			return nil, err
		}

		err = <-s.db.Message().Save(message)
		if err != nil {
			logger.Error("Failed to save message info", zap.Error(err))
			return nil, err
		}

		return &notificationPb.StatusResponse{
			Status: "success",
		}, nil

	}
	// Send message to the destination
	transactionId, err := s.messagingClient.SendMessage(req.TemplateId, req.RecipientPhoneNumber, parameters)
	if err != nil {
		return nil, err
	}

	// Save the message in the database
	message.TransactionId = transactionId
	err = <-s.db.Message().Save(message)
	if err != nil {
		logger.Error("Failed to save message info", zap.Error(err))
		return nil, err
	}

	return &notificationPb.StatusResponse{
		Status: "success",
	}, nil
}

func (s *MessagingService) RegisterMessagingTemplate(ctx context.Context, req *notificationPb.MessagingTemplate) (*notificationPb.StatusResponse, error) {
	userId, tenant := auth.GetUserIdAndTenant(ctx)

	// check if user is admin
	if !<-s.authClient.IsUserAdmin(ctx, userId) {
		logger.Error("User is not admin", zap.String("userId", userId))
		return nil, status.Error(codes.PermissionDenied, "User is not admin")
	}

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

func (s *MessagingService) FetchMessagingTemplates(ctx context.Context, req *notificationPb.FetchTemplateRequest) (*notificationPb.MessagingTemplateList, error) {
	userId, tenant := auth.GetUserIdAndTenant(ctx)

	// check if user is admin
	if !<-s.authClient.IsUserAdmin(ctx, userId) {
		logger.Error("User is not admin", zap.String("userId", userId))
		return nil, status.Error(codes.PermissionDenied, "User is not admin")
	}

	// Fetch the templates
	templates, totalCount := s.db.MessagingTemplate(tenant).GetTemplate(req)

	// Convert to protobuf
	templatesProto := make([]*notificationPb.MessagingTemplate, 0)
	for _, template := range templates {
		templatesProto = append(templatesProto, getMessagingTemplatePb(&template))
	}

	return &notificationPb.MessagingTemplateList{
		Templates:  templatesProto,
		TotalCount: int32(totalCount),
	}, nil
}

func (s *MessagingService) FetchMessages(ctx context.Context, req *notificationPb.FetchMessageRequest) (*notificationPb.MessageList, error) {
	userId, _ := auth.GetUserIdAndTenant(ctx)

	// check if user is admin
	if !<-s.authClient.IsUserAdmin(ctx, userId) {
		logger.Error("User is not admin", zap.String("userId", userId))
		return nil, status.Error(codes.PermissionDenied, "User is not admin")
	}

	// Fetch the messages
	messages, totalCount := s.db.Message().GetMessages(req.Filters, int64(req.PageNumber), int64(req.PageSize))

	// Convert to protobuf
	messagesProto := make([]*notificationPb.MessageProto, 0)
	for _, message := range messages {
		messagesProto = append(messagesProto, getMessageProto(&message))
	}

	return &notificationPb.MessageList{
		Messages:   messagesProto,
		TotalCount: int32(totalCount),
	}, nil
}

func (s *MessagingService) GetMessageMediaUploadUrl(ctx context.Context, req *notificationPb.MediaUploadRequest) (*notificationPb.MediaUploadUrl, error) {
	userId, tenant := auth.GetUserIdAndTenant(ctx)

	imagePath := fmt.Sprintf("whatsapp/%s/%s/%d/%d-image.%s", tenant, userId, time.Now().Year(), time.Now().Unix(), req.MediaExtension)
	socialBucket := os.Getenv("social_bucket")
	if socialBucket == "" {
		return nil, status.Error(codes.Internal, "social_bucket is not set")
	}

	acceptableExtensions := []string{"jpg", "jpeg", "png", "mp4", "webp", "doc", "pdf", "docx"}
	if !slices.Contains(acceptableExtensions, req.MediaExtension) {
		return nil, status.Error(codes.InvalidArgument, "Invalid media extension")
	}

	if req.MediaExtension == "" {
		req.MediaExtension = "jpg"
	}

	var contentType string

	if req.MediaExtension == "mp4" || req.MediaExtension == "webp" {
		contentType = fmt.Sprintf("video/%s", req.MediaExtension)
	} else if req.MediaExtension == "doc" || req.MediaExtension == "pdf" || req.MediaExtension == "docx" {
		contentType = fmt.Sprintf("application/%s", req.MediaExtension)
	} else {
		contentType = fmt.Sprintf("image/%s", req.MediaExtension)
	}

	uploadUrl, downloadUrl := s.cloudFns.GetPresignedUrl(socialBucket, imagePath, contentType, 15*time.Minute)
	return &notificationPb.MediaUploadUrl{
		UploadUrl: uploadUrl,
		MediaUrl:  downloadUrl,
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

func getMessagingTemplatePb(templateModel *models.MessagingTemplateModel) *notificationPb.MessagingTemplate {
	templateProto := &notificationPb.MessagingTemplate{}
	// Copy basic fields with copier
	copier.CopyWithOption(templateProto, templateModel, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// copy category
	templateProto.Category = notificationPb.Category(notificationPb.Category_value[templateModel.Category])

	// copy media parameters
	if templateModel.MediaParameters.MediaType != "" {
		templateProto.MediaParameters = &notificationPb.MediaParameters{
			MediaType: notificationPb.MediaType(notificationPb.MediaType_value[templateModel.MediaParameters.MediaType]),
			Link:      templateModel.MediaParameters.Link,
			Filename:  templateModel.MediaParameters.Filename,
		}
	}

	// copy button type
	templateProto.ButtonType = notificationPb.ButtonType(notificationPb.ButtonType_value[templateModel.ButtonType])

	// copy buttons
	if len(templateModel.Buttons.CallToActionButtons) > 0 {
		buttons := &notificationPb.Button{}
		for _, button := range templateModel.Buttons.CallToActionButtons {
			buttons.CallToActionButtons = append(buttons.CallToActionButtons, &notificationPb.CallToActionButtons{
				ActionType:  notificationPb.ActionType(notificationPb.ActionType_value[button.ActionType]),
				Text:        button.Text,
				PhoneNumber: button.PhoneNumber,
				Url: &notificationPb.Url{
					UrlType:       notificationPb.UrlType(notificationPb.UrlType_value[button.Url.UrlType]),
					Link:          button.Url.Link,
					UrlParameters: button.Url.UrlParameters,
				},
			})
		}
		templateProto.Buttons = buttons
	}

	return templateProto
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

func getMessageModel(sender string, req *notificationPb.MesssageRequest) *models.MessageModel {
	message := &models.MessageModel{
		Sender:     sender,
		Recipients: req.RecipientPhoneNumber,
		Message:    req.Preview,
	}

	if req.ScheduleInfo != nil {
		message.ScheduleInfo = models.ScheduleInfo{
			IsScheduled:   req.ScheduleInfo.IsScheduled,
			ScheduledTime: req.ScheduleInfo.ScheduledTime,
		}
	}

	if req.MediaParameters != nil {
		message.MediaParameters = models.MediaParameters{
			MediaType: req.MediaParameters.MediaType.String(),
			Link:      req.MediaParameters.Link,
			Filename:  req.MediaParameters.Filename,
		}
	}

	message.ButtonParameters = req.ButtonParameters

	return message
}

func getMessageProto(messageModel *models.MessageModel) *notificationPb.MessageProto {
	messageProto := &notificationPb.MessageProto{}
	copier.CopyWithOption(messageProto, messageModel, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return messageProto
}
