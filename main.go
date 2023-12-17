package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/Kotlang/notificationGo/generated"
	"github.com/Kotlang/notificationGo/jobs"
	"github.com/SaiNageswarS/go-api-boot/server"
)

var grpcPort = ":50051"
var webPort = ":8081"

func main() {
	// go-api-boot picks up keyvault name from environment variable.
	os.Setenv("AZURE-KEYVAULT-NAME", "kotlang-secrets")
	server.LoadSecretsIntoEnv(true)
	inject := NewInject()

	bootServer := server.NewGoApiBoot()
	pb.RegisterNotificationServiceServer(bootServer.GrpcServer, inject.NotificationService)

	// Jobs
	postCreated := jobs.NewPostCreatedJob(inject.NotificationDb)
	userCreated := jobs.NewUserCreatedJob(inject.NotificationDb)
	eventReminder := jobs.NewEventSubscribedJob(inject.NotificationDb)
	userFollow := jobs.NewUserFollowJob(inject.NotificationDb)
	inject.JobManager.RegisterJob(postCreated.Name, time.Minute*5, postCreated)
	inject.JobManager.RegisterJob(userCreated.Name, time.Minute*5, userCreated)
	inject.JobManager.RegisterJob(eventReminder.Name, time.Minute*1, eventReminder)
	inject.JobManager.RegisterJob(userFollow.Name, time.Minute*1, userFollow)
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
