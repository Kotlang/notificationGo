package models

import "github.com/Microsoft/go-winio/pkg/guid"

type EventModel struct {
	EventId            string            `bson:"_id" json:"eventId"`
	CreatorId          string            `bson:"creatorId" json:"creatorId"`
	EventType          string            `bson:"eventType" json:"eventType"`
	TemplateParameters map[string]string `bson:"templateParameters" json:"templateParameters"`
	IsBroadcast        bool              `bson:"isBroadcast" json:"isBroadcast"`
	TargetUsers        []string          `bson:"targetUsers" json:"targetUsers"`
}

func (m *EventModel) Id() string {
	if len(m.EventId) == 0 {
		g, _ := guid.NewV4()
		m.EventId = g.String()
	}
	return m.EventId
}
