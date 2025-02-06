package types

type LogCategoryCode string
type LogLevelCode string

const (
	LevelDebug LogLevelCode = "debug"
	LevelInfo  LogLevelCode = "info"
	LevelWarn  LogLevelCode = "warn"
	LevelError LogLevelCode = "error"
)

// Пример конвертации LogCategoryCode в строку
func (c LogCategoryCode) String() string {
	return string(c)
}
