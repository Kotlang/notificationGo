package main

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/service"
	"github.com/SaiNageswarS/go-api-boot/cloud"
	"github.com/SaiNageswarS/go-api-boot/jobs"
)

type Inject struct {
	NotificationDb db.NotificationDbInterface
	CloudFns       cloud.Cloud

	NotificationService *service.NotificationService

	JobManager *jobs.JobManager
}

func NewInject() *Inject {
	inj := &Inject{}

	inj.JobManager = jobs.NewJobManager("kotlang_jobs")

	inj.NotificationDb = &db.NotificationDb{}
	inj.CloudFns = &cloud.GCP{}

	inj.NotificationService = service.NewNotificationService(inj.NotificationDb)

	return inj
}
