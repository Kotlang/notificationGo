package jobs

import (
	"time"

	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type postCreated struct {
	Name string
	db   *db.NotificationDb
}

var Time = time.Now()

func NewPostCreatedJob(db *db.NotificationDb) *postCreated {
	return &postCreated{
		Name: "post.created",
		db:   db,
	}
}

func (j *postCreated) Run() (err error) {
	events := j.db.Event().GetEventByEventType(j.Name, 10, 0)

	if len(events) == 0 {
		return
	}

	for _, event := range events {
		title := event.TemplateParameters["title"]
		body := event.TemplateParameters["body"]

		err = extensions.SendMessageToTopic(title, body, event.Topic)
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
