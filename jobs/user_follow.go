package jobs

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userFollow struct {
	Name string
	db   *db.NotificationDb
}

func NewUserFollowJob(db *db.NotificationDb) *userFollow {
	return &userFollow{
		Name: "user.follow",
		db:   db,
	}
}

func (j *userFollow) Run() (err error) {
	events := j.db.Event().GetEventByEventType(j.Name, 10, 0)

	if len(events) == 0 {
		return
	}

	for _, event := range events {

		// if there are no target users, delete the event
		if len(event.TargetUsers) == 0 {
			logger.Error("Failed sending message to user", zap.Error(status.Error(codes.InvalidArgument, "no target users")))
			<-j.db.Event().DeleteById(event.Id())
			continue
		}

		// get device instance of followed user if err log the event and delete it so it doesn't block the queue
		DeviceInstance, devInstanceError := j.db.DeviceInstance().GetDeviceInstanceByUserId(event.TargetUsers[0])
		if devInstanceError != nil {
			logger.Error("Failed getting device instance", zap.Error(devInstanceError))
			<-j.db.Event().DeleteById(event.Id())
			continue
		}

		// send message to followed user if err log the event and delete it so it doesn't block the queue
		err = extensions.SendMessageToToken(event.Title, event.Body, event.ImageURL, DeviceInstance.Token, event.TemplateParameters)
		if err != nil {
			logger.Error("Failed sending message to user", zap.Error(err))
		}
		err = <-j.db.Event().DeleteById(event.Id())
		if err != nil {
			logger.Error("Failed deleting event", zap.Error(err))
			return
		}
	}

	return err
}
