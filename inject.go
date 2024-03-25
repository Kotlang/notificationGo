package main

import (
	"net/http"

	"github.com/Kotlang/notificationGo/clients"
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
	"github.com/Kotlang/notificationGo/service"
	"github.com/SaiNageswarS/go-api-boot/cloud"
	"github.com/SaiNageswarS/go-api-boot/jobs"
)

type Inject struct {
	NotificationDb db.NotificationDbInterface
	CloudFns       cloud.Cloud

	NotificationService *service.NotificationService
	MessaginService     *service.MessagingService
	Handlers            map[string]func(http.ResponseWriter, *http.Request)

	JobManager      *jobs.JobManager
	MessagingClient clients.MessagingClientInterface
}

func NewInject() *Inject {
	inj := &Inject{}

	inj.JobManager = jobs.NewJobManager("kotlang_jobs")

	inj.NotificationDb = &db.NotificationDb{}
	inj.CloudFns = &cloud.GCP{}
	inj.CloudFns.LoadSecretsIntoEnv()
	inj.MessagingClient = clients.NewWhatsAppClient()

	inj.NotificationService = service.NewNotificationService(inj.NotificationDb)
	inj.MessaginService = service.NewMessagingService(inj.MessagingClient, inj.NotificationDb)
	inj.Handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/messaging/incoming-message": extensions.WhatsappIncomingMessageHandler,
	}

	return inj
}
