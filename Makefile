BUILD_DIR=bin
APP=medgebot

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${APP} main.go
