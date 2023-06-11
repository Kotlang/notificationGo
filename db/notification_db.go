package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/odm"
)

type NotificationDb struct{}

func (a *NotificationDb) DeviceInstance() *DeviceInstanceRepository {
	baseRepo := odm.AbstractRepository[models.DeviceInstanceModel]{
		Database:       "kotlang_notification",
		CollectionName: "device_instance",
	}

	return &DeviceInstanceRepository{baseRepo}
}

func (a *NotificationDb) Event() *EventRepository {
	baseRepo := odm.AbstractRepository[models.EventModel]{
		Database:       "kotlang_notification",
		CollectionName: "event",
	}

	return &EventRepository{baseRepo}
}
