package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/auth"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Kotlang/notificationGo/generated"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	db *db.NotificationDb
}

var topics = []string{"post.created", "user.created", "event.created"}

func NewNotificationService(db *db.NotificationDb) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

func (s *NotificationService) RegisterDeviceInstance(ctx context.Context, req *pb.RegisterDeviceInstanceRequest) (*pb.NotificationStatusResponse, error) {
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
		return &pb.NotificationStatusResponse{
			Status: "success",
		}, nil
	}
}

func (s *NotificationService) RegisterEvent(ctx context.Context, req *pb.RegisterEventRequest) (*pb.NotificationStatusResponse, error) {
	creatorId, tenant := auth.GetUserIdAndTenant(ctx)

	if len(strings.TrimSpace(creatorId)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Token is not present.")
	}

	if len(strings.TrimSpace(req.EventType)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Event type is empty.")
	}

	topic := strings.TrimSpace(req.Topic)
	topicSplit := strings.Split(topic, ".")
	if (len(topicSplit) < 2) || (topicSplit[0] != tenant) {
		return nil, status.Error(codes.InvalidArgument, "Topic is invalid.")
	}

	event := &models.EventModel{
		CreatorId:          creatorId,
		EventType:          req.EventType,
		TemplateParameters: req.TemplateParameters,
		Topic:              topic,
		TargetUsers:        req.TargetUsers,
	}

	err := <-s.db.Event().Save(event)
	if err != nil {
		return nil, err
	} else {
		return &pb.NotificationStatusResponse{
			Status: "success",
		}, nil
	}
}

func (s *NotificationService) GetFCMTopics(ctx context.Context, req *pb.GetFCMTopicsRequest) (*pb.FCMTopicsResponse, error) {
	_, tenant := auth.GetUserIdAndTenant(ctx)

	return &pb.FCMTopicsResponse{
		Topics: funk.Map(topics, func(topic string) string { return fmt.Sprintf("%s.%s", tenant, topic) }).([]string),
	}, nil
}
