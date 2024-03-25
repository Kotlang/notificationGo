package models

import (
	"github.com/google/uuid"
)

type EventModel struct {
	EventId            string            `bson:"_id" json:"eventId"`
	CreatorId          string            `bson:"creatorId" json:"creatorId"`
	EventType          string            `bson:"eventType" json:"eventType"`
	Title              string            `bson:"title" json:"title"`
	Body               string            `bson:"body" json:"body"`
	ImageURL           string            `bson:"imageURL" json:"imageURL"`
	TemplateParameters map[string]string `bson:"templateParameters" json:"templateParameters"`
	Topic              string            `bson:"topic" json:"topic"`
	TargetUsers        []string          `bson:"targetUsers" json:"targetUsers"`
	Tenant             string            `bson:"tenant" json:"tenant"`
}

func (m *EventModel) Id() string {
	if len(m.EventId) == 0 {
		m.EventId = uuid.NewString()
	}
	return m.EventId
}
