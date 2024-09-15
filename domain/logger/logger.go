package logger

type LogLevel string

const (
	Info     LogLevel = "info"
	Error    LogLevel = "error"
	Critical LogLevel = "critical"
	Debug    LogLevel = "debug"
)

type Logger interface {
	Level() LogLevel
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Critical(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}
