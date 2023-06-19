package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type DeviceInstanceRepository struct {
	odm.AbstractRepository[models.DeviceInstanceModel]
}

func (e *DeviceInstanceRepository) GetDeviceInstance(pageNumber, pageSize int64) []models.DeviceInstanceModel {
	skip := pageNumber * pageSize
	resultChan, errChan := e.Find(bson.M{}, bson.D{}, pageSize, skip)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting device instance", zap.Error(err))
		return []models.DeviceInstanceModel{}
	}
}

func (e *DeviceInstanceRepository) CountFilteredDeviceInstance(redundantIds []string, tenant string) int64 {
	resultChan, errChan := e.CountDocuments(bson.M{"_id": bson.M{"$nin": redundantIds}, "tenant": tenant})

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed counting device instances without redundant device instances", zap.Error(err))
		return 0
	}
}

func (e *DeviceInstanceRepository) GetFilteredDeviceInstance(redundantIds []string, tenant string, pageNumber, pageSize int64) []models.DeviceInstanceModel {
	skip := pageNumber * pageSize
	resultChan, errChan := e.Find(bson.M{"_id": bson.M{"$nin": redundantIds}, "tenant": tenant}, bson.D{}, pageSize, skip)

	select {
	case res := <-resultChan:
		return res
	case err := <-errChan:
		logger.Error("Failed getting device instances without redundant device instances", zap.Error(err))
		return []models.DeviceInstanceModel{}
	}
}
