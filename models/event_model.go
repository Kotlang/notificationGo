package models

type EventModel struct {
	EventId            string            `bson:"_id" json:"eventId"`
	UserId             string            `bson:"userId" json:"userId"`
	EventType          string            `bson:"eventType" json:"eventType"`
	TemplateParameters map[string]string `bson:"templateParameters" json:"templateParameters"`
	IsBroadcast        bool              `bson:"isBroadcast" json:"isBroadcast"`
	TargetUsers        []string          `bson:"targetUsers" json:"targetUsers"`
}

func (m *EventModel) Id() string {
	return m.EventId
}
