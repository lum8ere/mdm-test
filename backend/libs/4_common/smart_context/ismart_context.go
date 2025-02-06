package smart_context

import (
	"context"
	"mdm/libs/4_common/types"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ISmartContext interface {
	// Методы логирования
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	// Методы управления полями логирования
	LogField(key string, value interface{}) ISmartContext  // будут заполнять поля в логгере
	LogFields(fields types.Fields) ISmartContext           // будут заполнять поля в логгере
	WithField(key string, value interface{}) ISmartContext // будут заполнять поля в DataField
	WithFields(fields types.Fields) ISmartContext          // будут заполнять поля в DataField

	GetLogger() *zap.Logger

	WithDbManager(db IDbManager) ISmartContext
	GetDbManager() IDbManager

	WithDB(db *gorm.DB) ISmartContext
	GetDB() *gorm.DB

	// Метод для получения стандартного context.Context
	WithContext(ctx context.Context) ISmartContext
	GetContext() context.Context
}
