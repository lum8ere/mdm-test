package handlers

import (
	"fmt"
	"strconv"
	"time"

	"mdm/libs/1_domain_methods/repositories"
	"mdm/libs/4_common/auth"
	"mdm/libs/4_common/smart_context"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Handler содержит зависимости для работы с устройствами.
type Handler struct {
	deviceRepo repositories.DeviceRepository
	userRepo   repositories.UserRepository
}

// NewHandler создаёт новый экземпляр Handler.
func NewHandler(repo repositories.DeviceRepository, userRepo repositories.UserRepository) *Handler {
	return &Handler{
		deviceRepo: repo,
		userRepo:   userRepo,
	}
}

// RegisterDeviceHandler обрабатывает регистрацию нового устройства.
// Ожидается, что в данных будет параметр "device_id".
func (h *Handler) RegisterDeviceHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	deviceID, ok := data["device_id"].(string)
	if !ok || deviceID == "" {
		return nil, fmt.Errorf("device_id is required")
	}
	return h.deviceRepo.RegisterDevice(sctx, deviceID)
}

// UpdateHeartbeatHandler обновляет время последнего обновления (heartbeat).
// Ожидается, что в данных будет параметр "id".
func (h *Handler) UpdateHeartbeatHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return h.deviceRepo.UpdateHeartbeat(sctx, id)
}

// GetDeviceStatusHandler возвращает статус устройства.
// Ожидается, что в данных будет параметр "id".
func (h *Handler) GetDeviceStatusHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return h.deviceRepo.GetDevice(sctx, id)
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

	return h.deviceRepo.SetCameraState(sctx, id, enabled)
}

func (h *Handler) UpdateMicrophoneHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
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
	return h.deviceRepo.SetMicrophoneState(sctx, id, enabled)
}

func (h *Handler) UpdateBluetoothHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
	enabled, ok := data["enabled"].(bool)
	if !ok {
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
	return h.deviceRepo.SetBluetoothState(sctx, id, enabled)
}

func (h *Handler) UpdateOsVersionHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
	version, ok := data["os_version"].(string)
	if !ok || version == "" {
		return nil, fmt.Errorf("os_version is required")
	}
	return h.deviceRepo.UpdateOsVersion(sctx, id, version)
}

func (h *Handler) UpdateBatteryLevelHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	id, ok := data["id"].(string)
	if !ok || id == "" {
		return nil, fmt.Errorf("id is required")
	}
	// В JSON числа обычно декодируются как float64
	levelVal, ok := data["battery_level"].(float64)
	if !ok {
		return nil, fmt.Errorf("battery_level is required and must be a number")
	}
	return h.deviceRepo.UpdateBatteryLevel(sctx, id, int(levelVal))
}

// Ожидается JSON: { "username": "...", "password": "..." }
func (h *Handler) LoginHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	username, ok1 := data["username"].(string)
	password, ok2 := data["password"].(string)
	if !ok1 || !ok2 || username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	user, err := h.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Сравнение захешированного пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Генерируем JWT-токен с информацией о пользователе
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(auth.JWTSecret)
	if err != nil {
		return nil, err
	}
	return map[string]string{"token": tokenString}, nil
}

func (h *Handler) GetAllDevicesHandler(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
	return h.deviceRepo.GetAllDevices(sctx)
}
