package jobs

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
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
		title := event.TemplateParameters["title"]
		body := event.TemplateParameters["body"]

		// get device instance of followed user
		DeviceInstance, devInstanceError := j.db.DeviceInstance().GetDeviceInstanceByUserId(event.TargetUsers[0])
		if devInstanceError != nil {
			logger.Error("Failed getting device instance", zap.Error(devInstanceError))
			continue
		}

		err = extensions.SendMessageToToken(title, body, DeviceInstance.Token, nil)
		if err != nil {
			return
		}
		err = <-j.db.Event().DeleteById(event.Id())
		if err != nil {
			return
		}
	}

	return err
}
