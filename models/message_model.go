package models

type MessageModel struct {
	MessageId string `bson:"_id" json:"messageId"`
	Sender    string `bson:"sender" json:"sender"`
	Receiver  string `bson:"receiver" json:"receiver"`
	Message   string `bson:"message" json:"message"`
	Status    string `bson:"status" json:"status"`
	CreatedOn int64  `bson:"createdOn" json:"createdOn"`
}
