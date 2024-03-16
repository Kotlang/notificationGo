package clients

import (
	"context"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

// FCMClient encapsulates Firebase Cloud Messaging client functionalities.
type FCMClient struct {
	app             *firebase.App
	fcmClient       *messaging.Client
	appCreationLock sync.Mutex
	ctx             context.Context
}

var fcmClientInstance *FCMClient

// NewFCMClient returns a singleton instance of FCMClient.
func NewFCMClient(ctx context.Context) *FCMClient {
	if fcmClientInstance == nil {
		fcmClientInstance = &FCMClient{
			ctx: ctx,
		}
		if err := fcmClientInstance.InitializeFCM(); err != nil {
			logger.Error("Failed to initialize FCM client", zap.Error(err))
		}
	}
	return fcmClientInstance
}

// InitializeFCM initializes Firebase Cloud Messaging client.
func (fc *FCMClient) InitializeFCM() error {
	fc.appCreationLock.Lock()
	defer fc.appCreationLock.Unlock()

	if fc.app == nil {
		opt := option.WithCredentialsJSON([]byte(os.Getenv("FCM-TOKEN")))
		ctx := context.Background()

		app, err := firebase.NewApp(ctx, nil, opt)
		if err != nil {
			logger.Error("Failed to initialize Firebase app", zap.Error(err))
			return err
		}

		client, err := app.Messaging(ctx)
		if err != nil {
			logger.Error("Failed to create FCM client", zap.Error(err))
			return err
		}

		fc.app = app
		fc.fcmClient = client
		fc.ctx = ctx
	}

	return nil
}

// SendMessageToToken sends a message to a single device token.
func (fc *FCMClient) SendMessageToToken(title, body, imageURL, token string, data map[string]string) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURL,
		},
		Data: data,
	}

	response, err := fc.fcmClient.Send(fc.ctx, message)
	logger.Info("FCM response: ", zap.String("response", response))
	return err
}

// SendMessageToMultipleTokens sends a message to multiple device tokens.
func (fc *FCMClient) SendMessageToMultipleTokens(title, body, imageURL string, tokens []string) error {
	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURL,
		},
		Tokens: tokens,
	}

	response, err := fc.fcmClient.SendMulticast(fc.ctx, message)
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

// SendMessageToTopic sends a message to a topic.
func (fc *FCMClient) SendMessageToTopic(title, body, imageURL, topic string) error {
	message := &messaging.Message{
		Topic: topic,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURL,
		},
	}

	response, err := fc.fcmClient.Send(fc.ctx, message)
	logger.Info("FCM response: ", zap.String("response", response))
	return err
}
