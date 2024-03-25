package jobs

import (
	"github.com/Kotlang/notificationGo/db"
	"github.com/Kotlang/notificationGo/extensions"
)

type WhatsAppJob struct {
	Name          string
	WhatsappCache *extensions.Cache
}

func NewWhatsAppJob(db db.NotificationDbInterface) *WhatsAppJob {
	return &WhatsAppJob{
		Name:          "whatsapp.job",
		WhatsappCache: extensions.GetCache(db),
	}
}

func (j *WhatsAppJob) Run() (err error) {
	j.WhatsappCache.UpdateDB()
	return
}
