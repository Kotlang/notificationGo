package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/odm"
)

type NotificationDbInterface interface {
	DeviceInstance() DeviceInstanceRepositoryInterface
	Event() EventRepositoryInterface
}

type NotificationDb struct{}

func (a *NotificationDb) DeviceInstance() DeviceInstanceRepositoryInterface {
	baseRepo := odm.UnimplementedBootRepository[models.DeviceInstanceModel]{
		Database:       "kotlang_notification",
		CollectionName: "device_instance",
	}

	return &DeviceInstanceRepository{baseRepo}
}

func (a *NotificationDb) Event() EventRepositoryInterface {
	baseRepo := odm.UnimplementedBootRepository[models.EventModel]{
		Database:       "kotlang_notification",
		CollectionName: "event",
	}

	return &EventRepository{baseRepo}
}
