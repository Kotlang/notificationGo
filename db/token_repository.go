package db

import (
	"github.com/Kotlang/notificationGo/models"
	"github.com/SaiNageswarS/go-api-boot/logger"
	"github.com/SaiNageswarS/go-api-boot/odm"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type DeviceInstanceRepositoryInterface interface {
	odm.BootRepository[models.DeviceInstanceModel]
	GetDeviceInstance(pageNumber, pageSize int64) []models.DeviceInstanceModel
	CountFilteredDeviceInstance(redundantIds []string, tenant string) int64
	GetFilteredDeviceInstance(redundantIds []string, tenant string, pageNumber, pageSize int64) []models.DeviceInstanceModel
	GetDeviceInstanceByUserId(userId string) (*models.DeviceInstanceModel, error)
	BulkGetDeviceInstanceByUserIds(userId []string) ([]models.DeviceInstanceModel, error)
}

type DeviceInstanceRepository struct {
	odm.UnimplementedBootRepository[models.DeviceInstanceModel]
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

// returns device instance containing FCM token by user id
func (e *DeviceInstanceRepository) GetDeviceInstanceByUserId(userId string) (*models.DeviceInstanceModel, error) {
	resultChan, errChan := e.FindOneById(userId)

	select {
	case res := <-resultChan:
		return res, nil
	case err := <-errChan:
		logger.Error("Failed getting device instance", zap.Error(err))
		return nil, err
	}
}

// returns device instance containing FCM token by user id
func (e *DeviceInstanceRepository) BulkGetDeviceInstanceByUserIds(userId []string) ([]models.DeviceInstanceModel, error) {
	filter := bson.M{"_id": bson.M{"$in": userId}}

	resultChan, errChan := e.Find(filter, bson.D{}, 0, 0)
	select {
	case res := <-resultChan:
		return res, nil
	case err := <-errChan:
		logger.Error("Failed getting device instance", zap.Error(err))
		return nil, err
	}
}
