package jobs

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type actionsNotify struct {
	Name string
	db   *db.NotificationDb
}

func NewActionsNotifyJob(db *db.NotificationDb) *actionsNotify {
	return &actionsNotify{
		Name: "actions.notify",
		db:   db,
	}
}

func (j *actionsNotify) Run() (err error) {
	events := j.db.Event().GetEventByEventType(j.Name, 10, 0)

	if len(events) == 0 {
		return
	}

	for _, event := range events {

		// if no target users delete the event log it and continue
		if len(event.TargetUsers) == 0 {
			logger.Error("Failed sending message to user",
				zap.Error(status.Error(codes.InvalidArgument, "no target users for event")))

			<-j.db.Event().DeleteById(event.Id())
			continue
		}

		// fetch fcmToken from db if err log the event and delete it so it doesn't block the queue
		fcmToken, fcmErr := j.db.DeviceInstance().GetDeviceInstanceByUserId(event.TargetUsers[0])
		if fcmErr != nil {
			logger.Error("Failed getting device instance", zap.Error(err), zap.String("userId", event.TargetUsers[0]))
			<-j.db.Event().DeleteById(event.Id())
			continue
		}

		// send message to fcmToken if err log the event and delete it so it doesn't block the queue
		err = extensions.SendMessageToToken(event.Title, event.Body, event.ImageURL, fcmToken.Token, event.TemplateParameters)

		if err != nil {
			logger.Error("Failed sending message to user", zap.Error(err))
		}

		err = <-j.db.Event().DeleteById(event.Id())
		if err != nil {
			logger.Error("Failed deleting event", zap.Error(err))
			return err
		}
	}
	return err
}
