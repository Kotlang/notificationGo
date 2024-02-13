package jobs

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type postCreated struct {
	Name string
	db   db.NotificationDbInterface
}

func NewPostCreatedJob(db db.NotificationDbInterface) *postCreated {
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

		// send message to topic if err log the event and delete it so it doesn't block the queue
		err = extensions.SendMessageToTopic(event.Title, event.Body, event.ImageURL, event.Topic)
		if err != nil {
			logger.Error("Failed sending message to topic", zap.Error(err), zap.String("postId", event.TemplateParameters["postId"]))
		}
		err = <-j.db.Event().DeleteById(event.Id())
		if err != nil {
			return
		}
	}

	return err
}
