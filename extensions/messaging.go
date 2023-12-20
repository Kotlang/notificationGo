package extensions

import (
	"errors"
	"sync"

	"firebase.google.com/go/v4/messaging"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
)

var fcm_client *FirebaseNotificationClient = &FirebaseNotificationClient{}

// FCM client
type FirebaseNotificationClient struct {
	cached_fcm_client        *messaging.Client
	fcm_client_creation_lock sync.Mutex
}

func (fc *FirebaseNotificationClient) GetFCMClient() *messaging.Client {
	fc.fcm_client_creation_lock.Lock()
	defer fc.fcm_client_creation_lock.Unlock()

	if fc.cached_fcm_client == nil {
		app := firebase_app.getFirestoreApp()
		if app == nil {
			logger.Error("Firebase app is nil")
			return nil
		}

		fcmClient, err := app.Messaging(firebase_app.ctx)
		if err != nil {
			logger.Error("Failed to create FCM client", zap.Error(err))
		}
		fc.cached_fcm_client = fcmClient
	}

	return fc.cached_fcm_client
}

func SendMessageToToken(title, body, token string, data map[string]string) error {
	fcmClient := fcm_client.GetFCMClient()
	if fcmClient == nil {
		return errors.New("FCM client is nil")
	}

	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	response, err := fcmClient.Send(firebase_app.ctx, message)

	logger.Info("FCM response: ", zap.String("response", response))

	return err
}

func SendMessageToMultipleTokens(title, body string, tokens []string) error {
	fcmClient := fcm_client.GetFCMClient()
	if fcmClient == nil {
		return errors.New("FCM client is nil")
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: tokens,
	}

	response, err := fcmClient.SendMulticast(firebase_app.ctx, message)

	if response != nil {
		logger.Info("FCM response: ", zap.Int("success_count", response.SuccessCount), zap.Int("failure_count", response.FailureCount))

		if response.Responses != nil && len(response.Responses) > 0 {
			for _, resp := range response.Responses {
				if resp != nil && resp.Error != nil {
					logger.Error("Response error: ", zap.String("error", resp.Error.Error()))
				}
			}
		}
	}

	return err
}

func SendMessageToTopic(title, body, topic string) error {
	fcmClient := fcm_client.GetFCMClient()
	if fcmClient == nil {
		return errors.New("FCM client is nil")
	}

	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	response, err := fcmClient.Send(firebase_app.ctx, message)

	logger.Info("FCM response: ", zap.String("response", response))

	return err
}
