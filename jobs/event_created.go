package jobs

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type eventCreated struct {
	Name string
	db   db.NotificationDbInterface
}

func NewEventCreatedJob(db db.NotificationDbInterface) *eventCreated {
	return &eventCreated{
		Name: "event.created",
		db:   db,
	}
}

func (j *eventCreated) Run() (err error) {
	events := j.db.Event().GetEventByEventType(j.Name, 10, 0)

	if len(events) == 0 {
		return
	}

	for _, event := range events {

		err = extensions.SendMessageToTopic(event.Title, event.Body, event.ImageURL, event.Topic)
		if err != nil {
			logger.Error("Failed sending message to topic", zap.Error(err))
			return
		}
		err = <-j.db.Event().DeleteById(event.Id())
		if err != nil {
			return
		}
	}

	return err
}
