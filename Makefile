BUILD_DIR=bin
APP=medgebot

.PHONY: build
build:
	go build -o ${BUILD_DIR}/${APP} main.go
