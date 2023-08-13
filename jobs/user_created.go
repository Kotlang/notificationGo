package jobs

import (
	"strings"

	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
)

type userCreated struct {
	Name string
	db   *db.NotificationDb
}

func NewUserCreatedJob(db *db.NotificationDb) *userCreated {
	return &userCreated{
		Name: "user.created",
		db:   db,
	}
}

func (j *userCreated) Run() (err error) {
	events := j.db.Event().GetEventByEventType(j.Name, 10, 0)

	if len(events) == 0 {
		return
	}

	for _, event := range events {
		title := event.TemplateParameters["title"]
		body := event.TemplateParameters["body"]

		tenant := strings.TrimSpace(event.Tenant)
		if len(tenant) != 0 {
			ids := event.RedundantUsers
			count := j.db.DeviceInstance().CountFilteredDeviceInstance(ids, tenant)
			if count > 0 {
				deviceInstance := j.db.DeviceInstance().GetFilteredDeviceInstance(ids, tenant, 0, count)

				tokens := funk.Map(deviceInstance, func(deviceInstance models.DeviceInstanceModel) string {
					return deviceInstance.Token
				}).([]string)

				err = extensions.SendMessageToMultipleTokens(title, body, tokens)

				if err != nil {
					logger.Error("Failed sending message to topic", zap.Error(err))
					return
				}
				err = <-j.db.Event().DeleteById(event.Id())
				if err != nil {
					return
				}
			}
		}
	}

	return err
}
