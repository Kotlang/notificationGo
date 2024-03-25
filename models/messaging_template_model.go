package models

type Url struct {
	UrlType       string            `bson:"urlType" json:"urlType"`
	Link          string            `bson:"link" json:"link"`
	UrlParameters map[string]string `bson:"urlParameters" json:"urlParameters"`
}

type CallToActionButtons struct {
	ActionType  string `bson:"actionType" json:"actionType"`
	Text        string `bson:"text" json:"text"`
	PhoneNumber string `bson:"phoneNumber,omitempty" json:"phoneNumber,omitempty"`
	Url         Url    `bson:"url,omitempty" json:"url,omitempty"`
}

type QuickReplyButtons struct {
	Text string `bson:"text" json:"text"`
}

type Button struct {
	CallToActionButtons []CallToActionButtons `bson:"callToActionButtons,omitempty" json:"callToActionButtons,omitempty"`
	QuickReplyButtons   []QuickReplyButtons   `bson:"quickReplyButtons,omitempty" json:"quickReplyButtons,omitempty"`
}

type MediaParameters struct {
	MediaType string `bson:"mediaType" json:"mediaType"`
	Link      string `bson:"link" json:"link"`
	Filename  string `bson:"filename,omitempty" json:"filename,omitempty"`
}

type MessagingTemplateModel struct {
	TemplateId       string            `bson:"_id" json:"templateId"`
	TemplateName     string            `bson:"templateName" json:"templateName"`
	MediaParameters  MediaParameters   `bson:"mediaParameters" json:"mediaParameters"`
	Header           string            `bson:"header,omitempty" json:"header,omitempty"`
	HeaderParameters map[string]string `bson:"headerParameters,omitempty" json:"headerParameters,omitempty"`
	Body             string            `bson:"body,omitempty" json:"body,omitempty"`
	BodyParameters   map[string]string `bson:"bodyParameters,omitempty" json:"bodyParameters,omitempty"`
	Footer           string            `bson:"footer,omitempty" json:"footer,omitempty"`
	Category         string            `bson:"category" json:"category"`
	WabaId           string            `bson:"wabaId,omitempty" json:"wabaId,omitempty"`
	ButtonType       string            `bson:"buttonType" json:"buttonType"`
	Buttons          Button            `bson:"buttons,omitempty" json:"buttons,omitempty"`
	CreatedOn        int64             `bson:"createdOn" json:"createdOn"`
}

func (m *MessagingTemplateModel) Id() string {
	return m.TemplateId
}
