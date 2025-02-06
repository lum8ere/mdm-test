package main

import (
	"mdm/libs/1_domain_methods/device_repository"
	"mdm/libs/1_domain_methods/handlers"
	"mdm/libs/1_domain_methods/run_processor"
	"mdm/libs/3_infrastructure/db_manager"
	"mdm/libs/4_common/env_vars"
	"mdm/libs/4_common/smart_context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	env_vars.LoadEnvVars()
	os.Setenv("LOG_LEVEL", "debug")

	logger := smart_context.NewSmartContext()

	dbm, err := db_manager.NewDbManager(logger)
	if err != nil {
		logger.Fatalf("Error connecting to database: %v", err)
	}
	logger = logger.WithDbManager(dbm)
	logger = logger.WithDB(dbm.GetGORM())

	// Инициализация репозитория устройств
	deviceRepo := device_repository.NewDeviceRepository(logger.GetDB())
	// Создаем хендлеры
	h := handlers.NewHandler(deviceRepo)

	r := chi.NewRouter()

	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With", "X-Request-Id", "X-Session-Id", "Apikey", "X-Api-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Регистрируем маршруты, используя обёртку JSONResponseMiddleware.
	// Обратите внимание: для параметра пути используем "{id}" (чтобы middleware мог извлечь его).
	r.Post("/devices/register", run_processor.JSONResponseMiddleware(logger, h.RegisterDeviceHandler))
	r.Post("/devices/{id}/heartbeat", run_processor.JSONResponseMiddleware(logger, h.UpdateHeartbeatHandler))
	r.Get("/devices/{id}/status", run_processor.JSONResponseMiddleware(logger, h.GetDeviceStatusHandler))
	r.Post("/devices/{id}/camera", run_processor.JSONResponseMiddleware(logger, h.UpdateCameraHandler))

	r.Post("/devices/{id}/microphone", run_processor.JSONResponseMiddleware(logger, h.UpdateMicrophoneHandler))
	r.Post("/devices/{id}/bluetooth", run_processor.JSONResponseMiddleware(logger, h.UpdateBluetoothHandler))
	r.Post("/devices/{id}/os", run_processor.JSONResponseMiddleware(logger, h.UpdateOsVersionHandler))
	r.Post("/devices/{id}/battery", run_processor.JSONResponseMiddleware(logger, h.UpdateBatteryLevelHandler))

	logger.Info("Server listening on port 4000")
	err = http.ListenAndServe(":4000", r)
	logger.Fatal(err)
}
