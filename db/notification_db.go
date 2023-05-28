package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/odm"
)

type NotificationDb struct{}

func (a *NotificationDb) DeviceInstance(tenant string) *DeviceInstanceRepository {
	baseRepo := odm.AbstractRepository[models.DeviceInstanceModel]{
		Database:       tenant + "_notification",
		CollectionName: "device_instance",
	}

	return &DeviceInstanceRepository{baseRepo}
}

func (a *NotificationDb) Event(tenant string) *EventRepository {
	baseRepo := odm.AbstractRepository[models.EventModel]{
		Database:       tenant + "_notification",
		CollectionName: "event",
	}

	return &EventRepository{baseRepo}
}
