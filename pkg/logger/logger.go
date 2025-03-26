package logger

import (
	"log"
)

type Logger struct{}

func InitLogger(level string) *Logger {
	// TODO logrus
	log.Printf("Инициализация логгера с уровнем: %s", level)
	return &Logger{}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf("[FATAL] "+format, args...)
}
