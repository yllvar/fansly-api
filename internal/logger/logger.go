package logger

import "log"

// Logger defines the interface for logging
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// New creates a new console logger
func New() Logger {
	return &consoleLogger{}
}

type consoleLogger struct{}

func (l *consoleLogger) Debugf(format string, args ...interface{}) {
	logf("DEBUG", format, args...)
}

func (l *consoleLogger) Infof(format string, args ...interface{}) {
	logf("INFO ", format, args...)
}

func (l *consoleLogger) Warnf(format string, args ...interface{}) {
	logf("WARN ", format, args...)
}

func (l *consoleLogger) Errorf(format string, args ...interface{}) {
	logf("ERROR", format, args...)
}

func logf(level, format string, args ...interface{}) {
	log.Printf("["+level+"] "+format, args...)
}
