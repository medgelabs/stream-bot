package logger

import (
	log "github.com/sirupsen/logrus"
)

// Trace logs a trace-level message
func Trace(msg string) {
	log.Trace(msg)
}

// Info logs a general info-level message
func Info(msg string) {
	log.Info(msg)
}

// Warn logs an warning-level message
func Warn(msg string) {
	log.Warn(msg)
}

// Error logs an error-level message
func Error(msg string, err error) {
	log.WithField("error", err).Error(msg)
}

// Fatal logs a fatal-level message, attaches the error() that caused it,
// and causes an app exit(1)
func Fatal(msg string, err error) {
	log.WithField("error", err).Fatal(msg)
}
