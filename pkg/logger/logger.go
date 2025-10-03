// Package logger provides a simple leveled logger implementation
package logger

import (
    "fmt"
    "log"
    "os"
    "strings"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
    // DebugLevel logs everything
    DebugLevel LogLevel = iota
    // InfoLevel logs informational messages and above
    InfoLevel
    // WarnLevel logs warnings and errors
    WarnLevel
    // ErrorLevel logs only errors
    ErrorLevel
    // FatalLevel logs only fatal errors
    FatalLevel
)

// Logger defines the interface for the logger
type Logger interface {
    Debug(args ...interface{})
    Info(args ...interface{})
    Warn(args ...interface{})
    Error(args ...interface{})
    Fatal(args ...interface{})
    Debugf(format string, args ...interface{})
    Infof(format string, args ...interface{})
    Warnf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
    Fatalf(format string, args ...interface{})
}

type logger struct {
    log      *log.Logger
    minLevel LogLevel
}

// New creates a new logger instance with the specified minimum log level
func New(level string) Logger {
    var minLevel LogLevel = InfoLevel // default to Info level
    
    switch strings.ToLower(level) {
    case "debug":
        minLevel = DebugLevel
    case "info":
        minLevel = InfoLevel
    case "warn", "warning":
        minLevel = WarnLevel
    case "error":
        minLevel = ErrorLevel
    case "fatal":
        minLevel = FatalLevel
    }

    return &logger{
        log:      log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
        minLevel: minLevel,
    }
}

func (l *logger) Debug(args ...interface{}) {
    if l.minLevel <= DebugLevel {
        l.log.Output(2, "[DEBUG] "+fmt.Sprint(args...))
    }
}

func (l *logger) Info(args ...interface{}) {
    if l.minLevel <= InfoLevel {
        l.log.Output(2, "[INFO] "+fmt.Sprint(args...))
    }
}

func (l *logger) Warn(args ...interface{}) {
    if l.minLevel <= WarnLevel {
        l.log.Output(2, "[WARN] "+fmt.Sprint(args...))
    }
}

func (l *logger) Error(args ...interface{}) {
    if l.minLevel <= ErrorLevel {
        l.log.Output(2, "[ERROR] "+fmt.Sprint(args...))
    }
}

func (l *logger) Fatal(args ...interface{}) {
    l.log.Output(2, "[FATAL] "+fmt.Sprint(args...))
    os.Exit(1)
}

func (l *logger) Debugf(format string, args ...interface{}) {
    if l.minLevel <= DebugLevel {
        l.log.Output(2, fmt.Sprintf("[DEBUG] "+format, args...))
    }
}

func (l *logger) Infof(format string, args ...interface{}) {
    if l.minLevel <= InfoLevel {
        l.log.Output(2, fmt.Sprintf("[INFO] "+format, args...))
    }
}

func (l *logger) Warnf(format string, args ...interface{}) {
    if l.minLevel <= WarnLevel {
        l.log.Output(2, fmt.Sprintf("[WARN] "+format, args...))
    }
}

func (l *logger) Errorf(format string, args ...interface{}) {
    if l.minLevel <= ErrorLevel {
        l.log.Output(2, fmt.Sprintf("[ERROR] "+format, args...))
    }
}

func (l *logger) Fatalf(format string, args ...interface{}) {
    l.log.Output(2, fmt.Sprintf("[FATAL] "+format, args...))
    os.Exit(1)
}