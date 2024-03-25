package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	notificationPb "github.com/Kotlang/notificationGo/generated/notification"
	"github.com/Kotlang/notificationGo/jobs"
	"github.com/SaiNageswarS/go-api-boot/server"
	"github.com/rs/cors"
)

var grpcPort = ":50051"
var webPort = ":8081"

func main() {

	inject := NewInject()
	inject.CloudFns.LoadSecretsIntoEnv()

	corsConfig := cors.New(
		cors.Options{
			AllowedHeaders: []string{"*"},
		})

	bootServer := server.NewGoApiBoot(corsConfig, nil, nil, inject.Handlers)
	notificationPb.RegisterNotificationServiceServer(bootServer.GrpcServer, inject.NotificationService)
	notificationPb.RegisterMessagingServiceServer(bootServer.GrpcServer, inject.MessaginService)

	//Jobs
	postCreated := jobs.NewPostCreatedJob(inject.NotificationDb)
	userCreated := jobs.NewUserCreatedJob(inject.NotificationDb)
	eventCreated := jobs.NewEventCreatedJob(inject.NotificationDb)
	eventReminder := jobs.NewEventReminderJob(inject.NotificationDb)
	userFollow := jobs.NewUserFollowJob(inject.NotificationDb)
	actionsNotify := jobs.NewActionsNotifyJob(inject.NotificationDb)
	whatsappMessage := jobs.NewWhatsappMessageJob(inject.NotificationDb)
	whatsappJob := jobs.NewWhatsAppJob(inject.NotificationDb)
	inject.JobManager.RegisterJob(postCreated.Name, time.Minute*5, postCreated)
	inject.JobManager.RegisterJob(userCreated.Name, time.Minute*5, userCreated)
	inject.JobManager.RegisterJob(eventCreated.Name, time.Minute*5, eventCreated)
	inject.JobManager.RegisterJob(eventReminder.Name, time.Minute, eventReminder)
	inject.JobManager.RegisterJob(userFollow.Name, time.Minute, userFollow)
	inject.JobManager.RegisterJob(actionsNotify.Name, time.Minute, actionsNotify)
	inject.JobManager.RegisterJob(whatsappMessage.Name, time.Minute, whatsappMessage)
	inject.JobManager.RegisterJob(whatsappJob.Name, time.Minute, whatsappJob)
	inject.JobManager.Start()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		inject.JobManager.Stop()
		bootServer.Stop()
	}()

	// Start the server
	bootServer.Start(grpcPort, webPort)
}
