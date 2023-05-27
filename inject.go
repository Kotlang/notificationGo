package main

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/service"
)

type Inject struct {
	NotificationDb *db.NotificationDb

	NotificationService *service.NotificationService
}

func NewInject() *Inject {
	inj := &Inject{}

	inj.NotificationDb = &db.NotificationDb{}

	inj.NotificationService = service.NewNotificationService(inj.NotificationDb)

	return inj
}
