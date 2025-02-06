package env_vars

import (
	"mdm/libs/4_common/smart_context"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

// GetCurrentFolder returns the directory of the current file.
func GetCurrentFolder() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func LoadEnvVars() {
	sctx := smart_context.NewSmartContext()
	_, filename, _, ok := runtime.Caller(0) // Получить путь текущего файла
	if !ok {
		sctx.Fatal("Error finding current file path")
	}

	basePath := filepath.Dir(filename) // Директория текущего файла
	sctx.Infof("Current file path: %s", basePath)

	envFile := os.Getenv("ENV_PATH")
	if envFile == "" {
		// Устанавливаем путь по умолчанию к .env
		envFile = filepath.Join(basePath, "..", "..", "..", ".env")
		sctx.Infof("ENV_PATH is not set, using default path: %s", envFile)
	} else {
		sctx.Infof("Found ENV_PATH='%s', so loading environment variables from this file...", envFile)
	}

	// Загрузка файла .env
	envFullFilePath := filepath.Clean(envFile) // Очистка пути
	err := godotenv.Load(envFullFilePath)
	if err != nil {
		// Текущая рабочая директория
		cwd, _ := os.Getwd()

		sctx.Fatalf("Error loading .env file: %v. Current directory '%s'", err.Error(), cwd)
	} else {
		sctx.Infof("Loaded .env file '%s'", envFullFilePath)
	}
}

// Helper to get an environment variable as an integer with a default value
func GetEnvAsInt(sctx smart_context.ISmartContext, key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		sctx.Infof("Invalid value for %s: %s. Using default: %d", key, valueStr, defaultValue)
		return defaultValue
	}
	return value
}
