package jobs

import (
	"context"
	"fmt"

	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type eventSubscribed struct {
	Name string
	db   *db.NotificationDb
}

func NewEventSubscribedJob(db *db.NotificationDb) *eventSubscribed {
	return &eventSubscribed{
		Name: "event.subscribed",
		db:   db,
	}
}

func (j *eventSubscribed) Run() (err error) {
	events := j.db.Event().GetEventByEventType(j.Name, 10, 0)

	if len(events) == 0 {
		return
	}

	for _, event := range events {
		title := event.TemplateParameters["title"]
		body := event.TemplateParameters["body"]

		subscriberIdList := <-extensions.GetEventSubscribers(context.Background(), event.TemplateParameters["eventId"])
		fmt.Println("subscriberIdList", subscriberIdList, title, body)
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
