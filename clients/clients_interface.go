package clients

type MessagingClient interface {
	SendMessage(templateID string, destination []string, parameters map[string]interface{}) (string, error)
}
