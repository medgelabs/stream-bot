package logger

import (
	log "github.com/sirupsen/logrus"
)

func Info(msg string) {
	log.Info(msg)
}

func Fatal(msg string, err error) {
	log.WithField("error", err).Fatal(msg)
}
