package main

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/service"
	"github.com/SaiNageswarS/go-api-boot/jobs"
)

type Inject struct {
	NotificationDb *db.NotificationDb

	NotificationService *service.NotificationService

	JobManager *jobs.JobManager
}

func NewInject() *Inject {
	inj := &Inject{}

	inj.JobManager = jobs.NewJobManager("navachar_jobs")

	inj.NotificationDb = &db.NotificationDb{}

	inj.NotificationService = service.NewNotificationService(inj.NotificationDb)

	return inj
}
