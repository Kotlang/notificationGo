package models

import "github.com/google/uuid"

type ScheduleInfo struct {
	IsScheduled   bool  `bson:"isScheduled" json:"isScheduled"`
	ScheduledTime int64 `bson:"scheduledTime" json:"scheduledTime"`
}

type MessageModel struct {
	MessageId        string            `bson:"_id" json:"messageId"`
	Sender           string            `bson:"sender" json:"sender"`
	Recipients       []string          `bson:"recipients" json:"recipients"`
	Message          string            `bson:"message" json:"message"`
	RecievedBy       []string          `bson:"recievedBy" json:"recievedBy"`
	ReadBy           []string          `bson:"readBy" json:"readBy"`
	RespondedBy      []string          `bson:"respondedBy" json:"respondedBy"`
	CreatedOn        int64             `bson:"createdOn" json:"createdOn"`
	ScheduleInfo     ScheduleInfo      `bson:"scheduleInfo" json:"scheduleInfo"`
	MediaParameters  MediaParameters   `bson:"mediaParameters" json:"mediaParameters"`
	ButtonParameters map[string]string `bson:"buttons" json:"buttons"`
	TransactionId    string            `bson:"transactionId" json:"transactionId"`
}

func (m *MessageModel) Id() string {
	if m.MessageId == "" {
		m.MessageId = uuid.NewString()
	}
	return m.MessageId
}
