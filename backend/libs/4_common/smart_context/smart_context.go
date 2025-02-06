package smart_context

import (
	"context"
	"fmt"
	"mdm/libs/4_common/types"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

const (
	DB_KEY = "db"
)

var _ ISmartContext = (*SmartContext)(nil)

func NewSmartContext() ISmartContext {
	config := zap.NewProductionConfig()
	logLevel := getLogLevel()
	atomicLevel := zap.NewAtomicLevelAt(logLevel)
	config.Level = atomicLevel

	config.DisableCaller = true
	config.EncoderConfig = zap.NewProductionEncoderConfig()

	if getLogToConsole() {
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Adds color to the level output
		config.DisableStacktrace = true
	} else {
		config.Encoding = "json"
	}

	// Build the logger from the modified configuration
	pureLogger, err := config.Build()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	sugarLogger := pureLogger.Sugar()

	syncCacheMap := &sync.Map{}
	sc := createSmartContext(sugarLogger, syncCacheMap, nil, nil, context.Background(), atomicLevel)

	return sc
}

type SmartContext struct {
	logger       *zap.SugaredLogger
	logFields    types.Fields    // поля которые были добавлены в логгер для логгирвания в каждом сообщении
	dataFields   types.Fields    // поля которые НЕ были добавлены в логгер, а используются для передачи данных между функциями - когда не хочется черзе кучу функций тащить параметры
	ctx          context.Context // для прерываний а также для дополнительных значений
	syncCacheMap *sync.Map
	logLevel     zap.AtomicLevel
}

func createSmartContext(
	innnerSugarLogger *zap.SugaredLogger,
	syncCacheMap *sync.Map,
	fields types.Fields,
	dataFields types.Fields,
	ctx context.Context,
	logLevel zap.AtomicLevel,
) *SmartContext {
	newPc := &SmartContext{
		logger:       innnerSugarLogger, // Инициализация logger
		syncCacheMap: syncCacheMap,
		logFields:    fields,
		dataFields:   dataFields,
		ctx:          ctx,
		logLevel:     logLevel, // Инициализация logLevel
	}
	return newPc
}

func (sc *SmartContext) Debug(args ...interface{}) {
	sc.logger.Debug(args...)
}

func (sc *SmartContext) Debugf(format string, args ...interface{}) {
	sc.logger.Debugf(format, args...)
}

func (sc *SmartContext) Error(args ...interface{}) {
	sc.logger.Error(args...)
}

func (sc *SmartContext) Errorf(format string, args ...interface{}) {
	sc.logger.Errorf(format, args...)
}

func (sc *SmartContext) Fatal(args ...interface{}) {
	sc.logger.Fatal(args...)
	panic(fmt.Sprint("", args))
}

func (sc *SmartContext) Fatalf(format string, args ...interface{}) {
	sc.logger.Fatalf(format, args...)
}

func (sc *SmartContext) GetContext() context.Context {
	return sc.ctx
}

func (sc *SmartContext) GetLogger() *zap.Logger {
	return sc.logger.Desugar()
}

func (sc *SmartContext) Info(args ...interface{}) {
	sc.logger.Info(args...)
}

func (sc *SmartContext) Infof(format string, args ...interface{}) {
	sc.logger.Infof(format, args...)
}

func (sc *SmartContext) LogField(key string, value interface{}) ISmartContext {
	newFields := sc.logFields.WithField(key, value)

	newPc := createSmartContext(
		sc.logger,
		sc.syncCacheMap,
		newFields,
		sc.dataFields,
		sc.ctx,
		sc.logLevel)
	return newPc
}

func (sc *SmartContext) LogFields(fields types.Fields) ISmartContext {
	newFields := sc.logFields.WithFields(fields)
	newPc := createSmartContext(
		sc.logger,
		sc.syncCacheMap,
		newFields,
		sc.dataFields,
		sc.ctx,
		sc.logLevel)

	return newPc
}

func (sc *SmartContext) Warn(args ...interface{}) {
	sc.logger.Warn(args...)
}

func (sc *SmartContext) Warnf(format string, args ...interface{}) {
	sc.logger.Warnf(format, args...)
}

func (sc *SmartContext) WithContext(ctx context.Context) ISmartContext {
	newPc := createSmartContext(
		sc.logger,
		sc.syncCacheMap,
		sc.logFields,
		sc.dataFields,
		ctx,
		sc.logLevel)
	return newPc
}

func (sc *SmartContext) WithField(key string, value interface{}) ISmartContext {
	newFields := sc.dataFields.WithField(key, value)
	newPc := createSmartContext(
		sc.logger,
		sc.syncCacheMap,
		sc.logFields,
		newFields,
		sc.ctx,
		sc.logLevel)
	return newPc
}

func (sc *SmartContext) WithFields(fields types.Fields) ISmartContext {
	newFields := sc.dataFields.WithFields(fields)
	newPc := createSmartContext(
		sc.logger,
		sc.syncCacheMap,
		sc.logFields,
		newFields,
		sc.ctx,
		sc.logLevel)
	return newPc
}

const DB_MANAGER_KEY = "db_manager"

func (sc *SmartContext) GetDbManager() IDbManager {
	result, ok := types.GetFieldTypedValue[IDbManager](sc.dataFields, DB_MANAGER_KEY)
	if !ok {
		return nil
	}
	return result
}

func (sc *SmartContext) WithDbManager(dbm IDbManager) ISmartContext {
	return sc.WithField(DB_MANAGER_KEY, dbm)
}

func (sc *SmartContext) WithDB(db *gorm.DB) ISmartContext {
	return sc.WithField(DB_KEY, db)
}

func (sc *SmartContext) GetDB() *gorm.DB {
	tx, ok := types.GetFieldTypedValue[*gorm.DB](sc.dataFields, DB_KEY)
	if !ok {
		return nil
	}
	if tx == nil {
		return nil
	}

	tx = tx.Session(&gorm.Session{NewDB: true, Context: sc.GetContext()})
	return tx
}

func getLogLevel() zapcore.Level {
	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(logLevel) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.WarnLevel // Default logging level
	}
}

func getLogToConsole() bool {
	logToConsole := os.Getenv("LOG_TO_CONSOLE")
	return len(logToConsole) > 0 // by default we are in PROD mode
}
