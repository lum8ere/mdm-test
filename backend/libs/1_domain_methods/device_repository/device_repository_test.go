package device_repository

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Используем in-memory SQLite для тестирования
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory sqlite: %v", err)
	}
	// Явное создание таблицы device без дефолтных функций
	createTableSQL := `
        CREATE TABLE device (
            id TEXT PRIMARY KEY,
            device_id TEXT NOT NULL,
            camera_enabled BOOLEAN NOT NULL,
            last_heartbeat DATETIME,
            created_at DATETIME,
            updated_at DATETIME
        );
    `
	if err := db.Exec(createTableSQL).Error; err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	return db
}

func TestRegisterDevice(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)

	deviceID := "test-device"
	device, err := repo.RegisterDevice(deviceID)
	if err != nil {
		t.Fatalf("Expected no error on registration, got: %v", err)
	}
	if device.DeviceID != deviceID {
		t.Errorf("Expected device_id %s, got %s", deviceID, device.DeviceID)
	}
}

func TestDuplicateRegistration(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)

	deviceID := "test-device"
	_, err := repo.RegisterDevice(deviceID)
	if err != nil {
		t.Fatalf("Expected first registration to succeed, got: %v", err)
	}

	_, err = repo.RegisterDevice(deviceID)
	if err == nil {
		t.Errorf("Expected error on duplicate registration, got nil")
	}
}

func TestUpdateHeartbeat(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)

	deviceID := "test-device"
	device, err := repo.RegisterDevice(deviceID)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	originalTime := device.LastHeartbeat

	// Немного подождём, чтобы время изменилось
	time.Sleep(1 * time.Second)

	updated, err := repo.UpdateHeartbeat(deviceID)
	if err != nil {
		t.Fatalf("UpdateHeartbeat failed: %v", err)
	}
	if !updated.LastHeartbeat.After(originalTime) {
		t.Errorf("Expected heartbeat time to be updated, original: %v, updated: %v", originalTime, updated.LastHeartbeat)
	}
}

func TestSetCameraState(t *testing.T) {
	db := setupTestDB(t)
	repo := NewDeviceRepository(db)

	deviceID := "test-device"
	_, err := repo.RegisterDevice(deviceID)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	updated, err := repo.SetCameraState(deviceID, true)
	if err != nil {
		t.Fatalf("SetCameraState failed: %v", err)
	}
	if !updated.CameraEnabled {
		t.Errorf("Expected CameraEnabled true, got false")
	}
}
