package logger

import (
	log "github.com/sirupsen/logrus"
)

func Info(msg string) {
	log.Info(msg)
}

func Panic(msg string, err error) {
	log.WithField("error", err).Fatal(msg)
}
