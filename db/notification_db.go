package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/odm"
)

type NotificationDbInterface interface {
	DeviceInstance() DeviceInstanceRepositoryInterface
	Event() EventRepositoryInterface
	MessagingTemplate(tenant string) MessagingTemplateRepositoryInterface
	Message(tenant string) MessageRepositoryInterface
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

func (a *NotificationDb) MessagingTemplate(tenant string) MessagingTemplateRepositoryInterface {
	baseRepo := odm.UnimplementedBootRepository[models.MessagingTemplateModel]{
		Database:       tenant + "_notification",
		CollectionName: "messaging_template",
	}

	return &MessagingTemplateRepository{baseRepo}
}

func (a *NotificationDb) Message(tenant string) MessageRepositoryInterface {
	baseRepo := odm.UnimplementedBootRepository[models.MessageModel]{
		Database:       tenant + "_notification",
		CollectionName: "message",
	}

	return &MessageRepository{baseRepo}
}
