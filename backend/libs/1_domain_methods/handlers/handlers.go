package handlers

import (
	"fmt"
	"strconv"

	"mdm/libs/1_domain_methods/device_repository"
	"mdm/libs/4_common/smart_context"
)

// Handler содержит зависимости для работы с устройствами.
type Handler struct {
	deviceRepo device_repository.DeviceRepository
}

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(repo device_repository.DeviceRepository) *Handler {
	return &Handler{
		deviceRepo: repo,
	}
}

// RegisterDeviceHandler обрабатывает регистрацию нового устройства.
// Ожидается, что в данных будет параметр "device_id".
func (h *Handler) RegisterDeviceHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	deviceID, ok := data["device_id"].(string)
	if !ok || deviceID == "" {
		return nil, fmt.Errorf("device_id is required")
	}
	return h.deviceRepo.RegisterDevice(deviceID)
}

// UpdateHeartbeatHandler обновляет время последнего обновления (heartbeat).
// Ожидается, что в данных будет параметр "id".
func (h *Handler) UpdateHeartbeatHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return h.deviceRepo.UpdateHeartbeat(id)
}

// GetDeviceStatusHandler возвращает статус устройства.
// Ожидается, что в данных будет параметр "id".
func (h *Handler) GetDeviceStatusHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return h.deviceRepo.GetDevice(id)
}

// UpdateCameraHandler изменяет состояние камеры устройства.
// Ожидается, что в данных будет параметр "id" и "enabled".
func (h *Handler) UpdateCameraHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}

	// Попытка извлечь булево значение из поля "enabled"
	enabled, ok := data["enabled"].(bool)
	if !ok {
		// Если значение передано как строка, попробуем преобразовать
		if strVal, ok := data["enabled"].(string); ok {
			parsed, err := strconv.ParseBool(strVal)
			if err != nil {
				return nil, fmt.Errorf("invalid value for enabled")
			}
			enabled = parsed
		} else {
			return nil, fmt.Errorf("enabled parameter is required and must be boolean")
		}
	}

	return h.deviceRepo.SetCameraState(id, enabled)
}
