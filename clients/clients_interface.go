package clients

import "context"

// MessagingClient defines the interface for interacting with the messaging service.
type MessagingClientInterface interface {
	SendMessage(templateID string, destination []string, parameters map[string]interface{}) (string, error)
}

// NotificationClientInterface defines the interface for interacting with the Firebase notification client.
type NotificationClientInterface interface {
	// SendMessageToToken sends a message to a single device token.
	SendMessageToToken(title, body, imageURL, token string, data map[string]string) error
	// SendMessageToMultipleTokens sends a message to multiple device tokens.
	SendMessageToMultipleTokens(title, body, imageURL string, tokens []string) error
	// SendMessageToTopic sends a message to subscribers of a topic
	SendMessageToTopic(title, body, imageURL, topic string) error
}

// SocialInterface defines the interface for interacting with the social service.
type SocialClientInterface interface {
	// GetEventSubscribers retrieves subscribers for a given event asynchronously.
	GetEventSubscribers(ctx context.Context, tenant string, eventId string) chan []string
}

// AuthInterface defines the interface for interacting with the authentication service.
type AuthClientInterface interface {
	// IsUserAdmin checks if a user is an admin.
	IsUserAdmin(ctx context.Context, userId string) chan bool
}
