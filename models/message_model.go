package models

type MessageModel struct {
	MessageId   string   `bson:"_id" json:"messageId"`
	Sender      string   `bson:"sender" json:"sender"`
	Recipients  []string `bson:"recipients" json:"recipients"`
	Message     string   `bson:"message" json:"message"`
	RecievedBy  []string `bson:"recievedBy" json:"recievedBy"`
	ReadBy      []string `bson:"readBy" json:"readBy"`
	RespondedBy []string `bson:"respondedBy" json:"respondedBy"`
	CreatedOn   int64    `bson:"createdOn" json:"createdOn"`
}
