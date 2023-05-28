package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/odm"
)

type EventRepository struct {
	odm.AbstractRepository[models.EventModel]
}
