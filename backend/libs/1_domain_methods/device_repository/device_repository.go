package device_repository

import (
	"errors"
	"fmt"
	"mdm/libs/2_generated_models/model"
	"time"

	"gorm.io/gorm"
)

// DeviceRepository описывает набор операций над устройствами.
type DeviceRepository interface {
	RegisterDevice(deviceID string) (*model.Device, error)
	GetDevice(deviceID string) (*model.Device, error)
	UpdateHeartbeat(deviceID string) (*model.Device, error)
	SetCameraState(deviceID string, enabled bool) (*model.Device, error)
}

// repository — реализация DeviceRepository, использующая GORM.
type repository struct {
	db *gorm.DB
}

// NewDeviceRepository возвращает новый экземпляр репозитория.
func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &repository{db: db}
}

// RegisterDevice регистрирует устройство, если оно ещё не зарегистрировано.
func (r *repository) RegisterDevice(deviceID string) (*model.Device, error) {
	// Проверяем, существует ли уже устройство.
	var existing model.Device
	err := r.db.Where("device_id = ?", deviceID).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("device already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	device := &model.Device{
		DeviceID:      deviceID,
		CameraEnabled: false,
		LastHeartbeat: time.Now(),
	}
	if err := r.db.Create(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

// GetDevice возвращает данные об устройстве по его DeviceID.
func (r *repository) GetDevice(deviceID string) (*model.Device, error) {
	var device model.Device
	if err := r.db.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// UpdateHeartbeat обновляет время последней активности устройства.
func (r *repository) UpdateHeartbeat(deviceID string) (*model.Device, error) {
	device, err := r.GetDevice(deviceID)
	if err != nil {
		return nil, err
	}
	device.LastHeartbeat = time.Now()
	if err := r.db.Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

// SetCameraState изменяет состояние камеры у устройства.
func (r *repository) SetCameraState(deviceID string, enabled bool) (*model.Device, error) {
	device, err := r.GetDevice(deviceID)
	if err != nil {
		return nil, err
	}
	device.CameraEnabled = enabled
	if err := r.db.Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}
