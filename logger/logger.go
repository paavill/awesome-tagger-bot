package logger

import (
	"log"
	"os"

	"github.com/paavill/awesome-tagger-bot/domain/logger"
)

type lgr struct {
	logger   *log.Logger
	level    uint8
	rawLevel logger.LogLevel
}

func New(logLevel logger.LogLevel) logger.Logger {
	level := getLevel(logLevel)
	stdLogger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
	return &lgr{
		logger:   stdLogger,
		level:    level,
		rawLevel: logLevel,
	}
}

func (l *lgr) Level() logger.LogLevel {
	return l.rawLevel
}

func (l *lgr) Debug(msg string, args ...interface{}) {
	l.logger.SetPrefix("[DEBUG] ")
	if l.level >= 3 {
		l.logger.Printf(msg, args...)
	}
}

func (l *lgr) Info(msg string, args ...interface{}) {
	l.logger.SetPrefix("[INFO] ")
	if l.level >= 2 {
		l.logger.Printf(msg, args...)
	}
}

func (l *lgr) Error(msg string, args ...interface{}) {
	l.logger.SetPrefix("[ERROR] ")
	if l.level >= 1 {
		l.logger.Printf(msg, args...)
	}
}

func (l *lgr) Critical(msg string, args ...interface{}) {
	l.logger.SetPrefix("[CRITICAL] ")
	l.logger.Printf(msg, args...)
}

func getLevel(logLevel logger.LogLevel) uint8 {
	switch logLevel {
	case logger.Critical:
		return 0
	case logger.Error:
		return 1
	case logger.Info:
		return 2
	case logger.Debug:
		return 3
	default:
		return 2
	}
}
