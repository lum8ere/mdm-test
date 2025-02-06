package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Device соответствует JSON-структуре, возвращаемой сервером (см. модель Device в базе)
type Device struct {
	ID            string    `json:"id"`
	DeviceID      string    `json:"device_id"`
	CameraEnabled bool      `json:"camera_enabled"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// registerDevice отправляет запрос на регистрацию устройства (POST /devices/register)
func registerDevice(server, deviceID string) (*Device, error) {
	url := fmt.Sprintf("%s/devices/register", server)
	payload := map[string]string{
		"device_id": deviceID,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Если сервер вернул не OK, читаем тело ответа для диагностики
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("registration failed: %s", body)
	}

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}
	return &device, nil
}

// sendHeartbeat отправляет запрос heartbeat (POST /devices/{device_id}/heartbeat)
func sendHeartbeat(server, deviceID string) (*Device, error) {
	url := fmt.Sprintf("%s/devices/%s/heartbeat", server, deviceID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("heartbeat failed: %s", body)
	}

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}
	return &device, nil
}

// getDeviceStatus получает статус устройства (GET /devices/{device_id}/status)
func getDeviceStatus(server, deviceID string) (*Device, error) {
	url := fmt.Sprintf("%s/devices/%s/status", server, deviceID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("get status failed: %s", body)
	}

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}
	return &device, nil
}

func main() {
	// Парсинг флагов командной строки
	deviceID := flag.String("device-id", "", "Уникальный идентификатор устройства")
	serverURL := flag.String("server", "http://localhost:4000", "URL сервера MDM")
	flag.Parse()

	if *deviceID == "" {
		fmt.Println("Параметр --device-id обязателен")
		os.Exit(1)
	}

	// Регистрируем устройство
	device, err := registerDevice(*serverURL, *deviceID)
	if err != nil {
		log.Printf("Ошибка регистрации устройства: %v", err)
		// Если устройство уже зарегистрировано, попробуем получить его статус
		log.Printf("Пытаемся получить статус устройства")
		device, err = getDeviceStatus(*serverURL, *deviceID)
		if err != nil {
			log.Fatalf("Не удалось получить статус устройства: %v", err)
		}
	}
	log.Printf("Устройство зарегистрировано: %+v", device)

	// Сохраним текущее состояние камеры для отслеживания изменений
	currentCameraState := device.CameraEnabled

	// Периодически отправляем heartbeat (например, каждые 10 секунд)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		updatedDevice, err := sendHeartbeat(*serverURL, *deviceID)
		if err != nil {
			log.Printf("Ошибка отправки heartbeat: %v", err)
			continue
		}
		log.Printf("Получен heartbeat: %+v", updatedDevice)

		// Если изменилось состояние камеры, логируем событие
		if updatedDevice.CameraEnabled != currentCameraState {
			currentCameraState = updatedDevice.CameraEnabled
			if currentCameraState {
				log.Printf("Команда: Включить камеру для устройства %s", *deviceID)
			} else {
				log.Printf("Команда: Выключить камеру для устройства %s", *deviceID)
			}
		}
	}
}
