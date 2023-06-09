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
	createPost := jobs.NewCreatePostJob()
	inject.JobManager.RegisterJob(createPost.Name, time.Second*10, createPost)
	inject.JobManager.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		inject.JobManager.Stop()
		bootServer.Stop()
	}()

	pb.RegisterNotificationServiceServer(bootServer.GrpcServer, inject.NotificationService)

	bootServer.Start(grpcPort, webPort)
}
