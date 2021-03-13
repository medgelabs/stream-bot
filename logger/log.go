package logger

import (
	log "github.com/sirupsen/logrus"
)

// Trace logs a trace-level message
func Trace(format string, args ...interface{}) {
	log.Tracef(format, args...)
}

// Info logs a general info-level message
func Info(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warn logs an warning-level message
func Warn(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Error logs an error-level message
func Error(err error, format string, args ...interface{}) {
	log.WithField("error", err).Errorf(format, args...)
}

// Fatal logs a fatal-level message, attaches the error() that caused it,
// and causes an app exit(1)
func Fatal(err error, format string, args ...interface{}) {
	log.WithField("error", err).Fatalf(format, args...)
}
