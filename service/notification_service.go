package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/auth"
	"github.com/jinzhu/copier"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	notificationPb "github.com/Kotlang/notificationGo/generated/notification"
)

type NotificationService struct {
	notificationPb.UnimplementedNotificationServiceServer
	db db.NotificationDbInterface
}

var topics = []string{"post.created", "event.created"}

func NewNotificationService(db db.NotificationDbInterface) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

func (s *NotificationService) RegisterDeviceInstance(ctx context.Context, req *notificationPb.RegisterDeviceInstanceRequest) (*notificationPb.NotificationStatusResponse, error) {
	userId, tenant := auth.GetUserIdAndTenant(ctx)

	if len(strings.TrimSpace(userId)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Token is not present.")
	}

	deviceInstance := &models.DeviceInstanceModel{
		LoginId: userId,
		Token:   req.Token,
		Tenant:  tenant,
	}

	err := <-s.db.DeviceInstance().Save(deviceInstance)
	if err != nil {
		return nil, err
	} else {
		return &notificationPb.NotificationStatusResponse{
			Status: "success",
		}, nil
	}
}

func (s *NotificationService) RegisterEvent(ctx context.Context, req *notificationPb.RegisterEventRequest) (*notificationPb.NotificationStatusResponse, error) {
	creatorId, tenant := auth.GetUserIdAndTenant(ctx)

	if len(strings.TrimSpace(creatorId)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Token is not present.")
	}

	if len(strings.TrimSpace(req.EventType)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Event type is empty.")
	}

	// check if topic is valid and belongs to the tenant
	topic := strings.TrimSpace(req.Topic)
	topicSplit := strings.Split(topic, ".")
	if (len(topicSplit) < 2) || (topicSplit[0] != tenant) {
		return nil, status.Error(codes.InvalidArgument, "Topic is invalid.")
	}

	// copy request to event model
	event := &models.EventModel{}
	copier.CopyWithOption(event, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// save event to db and return response to client
	err := <-s.db.Event().Save(event)
	if err != nil {
		return nil, err
	} else {
		return &notificationPb.NotificationStatusResponse{
			Status: "success",
		}, nil
	}
}

func (s *NotificationService) GetFCMTopics(ctx context.Context, req *notificationPb.GetFCMTopicsRequest) (*notificationPb.FCMTopicsResponse, error) {
	_, tenant := auth.GetUserIdAndTenant(ctx)

	return &notificationPb.FCMTopicsResponse{
		Topics: funk.Map(topics, func(topic string) string { return fmt.Sprintf("%s.%s", tenant, topic) }).([]string),
	}, nil
}
