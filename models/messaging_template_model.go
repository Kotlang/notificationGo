package models

import (
	"github.com/google/uuid"
)

type MediaParameters struct {
	MediaType string `bson:"mediaType" json:"mediaType"`
	Link      string `bson:"link" json:"link"`
	Filename  string `bson:"filename" json:"filename"`
}

type MessagingTemplateModel struct {
	TemplateId         string            `bson:"_id" json:"templateId"`
	TemplateName       string            `bson:"templateName" json:"templateName"`
	TemplateParameters map[string]string `bson:"templateParameters" json:"templateParameters"`
	MediaParameters    MediaParameters   `bson:"mediaParameters" json:"mediaParameters"`
}

func (m *MessagingTemplateModel) Id() string {
	if len(m.TemplateId) == 0 {
		m.TemplateId = uuid.NewString()
	}
	return m.TemplateId
}
