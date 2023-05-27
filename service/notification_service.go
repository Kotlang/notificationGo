package service

import (
	"context"
	"strings"

	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/Kotlang/notificationGo/generated"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
	db *db.NotificationDb
}

func NewNotificationService(db *db.NotificationDb) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

func (s *NotificationService) RegisterDeviceInstance(ctx context.Context, req *pb.RegisterDeviceInstanceRequest) (*pb.StatusResponse, error) {
	userId, tenant := auth.GetUserIdAndTenant(ctx)

	if len(strings.TrimSpace(userId)) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Token is not present.")
	}

	deviceInstance := &models.DeviceInstanceModel{
		LoginId: userId,
		Token:   req.Token,
	}

	err := <-s.db.DeviceInstance(tenant).Save(deviceInstance)
	if err != nil {
		return nil, err
	} else {
		return &pb.StatusResponse{
			Status: "success",
		}, nil
	}
}
