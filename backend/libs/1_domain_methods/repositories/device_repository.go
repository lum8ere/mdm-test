package repositories

import (
	"errors"
	"mdm/libs/2_generated_models/model"
	"mdm/libs/4_common/smart_context"
	"time"

	"gorm.io/gorm"
)

// DeviceRepository описывает набор операций над устройствами.
type DeviceRepository interface {
	RegisterDevice(sctx smart_context.ISmartContext, deviceID string) (*model.Device, error)
	GetDevice(sctx smart_context.ISmartContext, deviceID string) (*model.Device, error)
	UpdateHeartbeat(sctx smart_context.ISmartContext, deviceID string) (*model.Device, error)
	SetCameraState(sctx smart_context.ISmartContext, deviceID string, enabled bool) (*model.Device, error)
	SetMicrophoneState(sctx smart_context.ISmartContext, deviceID string, enabled bool) (*model.Device, error)
	SetBluetoothState(sctx smart_context.ISmartContext, deviceID string, enabled bool) (*model.Device, error)
	UpdateOsVersion(sctx smart_context.ISmartContext, deviceID string, version string) (*model.Device, error)
	UpdateBatteryLevel(sctx smart_context.ISmartContext, deviceID string, level int) (*model.Device, error)
	GetAllDevices(sctx smart_context.ISmartContext) ([]model.Device, error)
}

// repository — реализация DeviceRepository, использующая GORM.
type device_repository struct {
	db *gorm.DB
}

// NewDeviceRepository возвращает новый экземпляр репозитория.
func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &device_repository{db: db}
}

// RegisterDevice регистрирует устройство, если оно ещё не зарегистрировано.
func (r *device_repository) RegisterDevice(sctx smart_context.ISmartContext, deviceID string) (*model.Device, error) {
	// Проверяем, существует ли уже устройство.
	var existing model.Device
	err := r.db.Where("device_id = ?", deviceID).First(&existing).Error
	if err == nil {
		sctx.Warnf("device already registered")
		return nil, nil
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
	sctx.Infof("device registered")
	return device, nil
}

// GetDevice возвращает данные об устройстве по его DeviceID.
func (r *device_repository) GetDevice(sctx smart_context.ISmartContext, deviceID string) (*model.Device, error) {
	var device model.Device
	if err := r.db.Where("device_id = ?", deviceID).First(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

// UpdateHeartbeat обновляет время последней активности устройства.
func (r *device_repository) UpdateHeartbeat(sctx smart_context.ISmartContext, deviceID string) (*model.Device, error) {
	device, err := r.GetDevice(sctx, deviceID)
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
func (r *device_repository) SetCameraState(sctx smart_context.ISmartContext, deviceID string, enabled bool) (*model.Device, error) {
	device, err := r.GetDevice(sctx, deviceID)
	if err != nil {
		return nil, err
	}
	device.CameraEnabled = enabled
	if err := r.db.Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

func (r *device_repository) SetMicrophoneState(sctx smart_context.ISmartContext, deviceID string, enabled bool) (*model.Device, error) {
	device, err := r.GetDevice(sctx, deviceID)
	if err != nil {
		return nil, err
	}
	device.MicrophoneEnabled = enabled
	if err := r.db.Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

func (r *device_repository) SetBluetoothState(sctx smart_context.ISmartContext, deviceID string, enabled bool) (*model.Device, error) {
	device, err := r.GetDevice(sctx, deviceID)
	if err != nil {
		return nil, err
	}
	device.BluetoothEnabled = enabled
	if err := r.db.Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

func (r *device_repository) UpdateOsVersion(sctx smart_context.ISmartContext, deviceID string, version string) (*model.Device, error) {
	device, err := r.GetDevice(sctx, deviceID)
	if err != nil {
		return nil, err
	}
	device.OsVersion = version
	if err := r.db.Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

func (r *device_repository) UpdateBatteryLevel(sctx smart_context.ISmartContext, deviceID string, level int) (*model.Device, error) {
	device, err := r.GetDevice(sctx, deviceID)
	if err != nil {
		return nil, err
	}
	device.BatteryLevel = int32(level)
	if err := r.db.Save(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

func (r *device_repository) GetAllDevices(sctx smart_context.ISmartContext) ([]model.Device, error) {
	var devices []model.Device
	if err := r.db.Find(&devices).Error; err != nil {
		return nil, err
	}

	sctx.Infof("len(devices) = %d", len(devices))
	return devices, nil
}
