package jobs

import (
	"context"
	"strconv"
	"time"

	"github.com/Kotlang/notificationGo/clients"
	"github.com/Kotlang/notificationGo/db"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

type eventReminder struct {
	Name               string
	db                 db.NotificationDbInterface
	notificationClient clients.NotificationClientInterface
	socialClient       clients.SocialClientInterface
}

func NewEventReminderJob(db db.NotificationDbInterface) *eventReminder {
	return &eventReminder{
		Name:               "event.reminder",
		db:                 db,
		notificationClient: clients.NewFCMClient(context.Background()),
		socialClient:       clients.NewSocialClient(),
	}
}

func (j *eventReminder) Run() (err error) {
	events := j.db.Event().GetEventByEventType(j.Name, 10, 0)

	if len(events) == 0 {
		return
	}

	for _, event := range events {

		// parse the event start time if it is not parsable, delete the event and log the eventId
		eventStartTime, intErr := strconv.ParseInt(event.TemplateParameters["startAt"], 10, 64)
		if intErr != nil {
			logger.Error("Failed parsing event start time", zap.Error(intErr), zap.String("event", event.TemplateParameters["eventId"]))
			err = <-j.db.Event().DeleteById(event.Id())
		}

		// if event start time is more than 10 minutes from now, skip
		if eventStartTime-time.Now().Unix() >= 600 {
			continue
		}

		// if event start time is less than 0 minutes from now, delete the event and log the eventId
		if eventStartTime-time.Now().Unix() < 0 {
			logger.Error("Event start time is less than current time", zap.String("event", event.TemplateParameters["eventId"]))
			err = <-j.db.Event().DeleteById(event.Id())
			continue
		}

		// if event start time is less than 10 minutes from now, send the notification
		subscriberIdList := <-j.socialClient.GetEventSubscribers(context.TODO(), event.Tenant, event.TemplateParameters["eventId"])

		// if there are no subscribers, delete the event and log the eventId
		if len(subscriberIdList) == 0 {
			logger.Error("No subscribers found for event", zap.String("event", event.TemplateParameters["eventId"]))
			err = <-j.db.Event().DeleteById(event.Id())
			continue
		}

		FCMTokenList, errRes := j.db.DeviceInstance().BulkGetDeviceInstanceByUserIds(subscriberIdList)
		if errRes != nil {
			logger.Error("Failed getting device instance", zap.Error(errRes))
			err = <-j.db.Event().DeleteById(event.Id())
			continue
		}

		fcmIds := make([]string, 0)
		for _, fcmToken := range FCMTokenList {
			fcmIds = append(fcmIds, fcmToken.Token)
		}
		err = j.notificationClient.SendMessageToMultipleTokens(event.Title, event.Body, event.ImageURL, fcmIds)
		if err != nil {
			logger.Error("Failed to send message", zap.Error(err), zap.String("event", event.TemplateParameters["eventId"]))
		}

		err = <-j.db.Event().DeleteById(event.Id())
		if err != nil {
			logger.Error("Failed to delete event", zap.Error(err))
			return
		}
	}
	return err
}
