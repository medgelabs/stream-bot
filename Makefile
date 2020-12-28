BUILD_DIR=bin
APP=medgebot

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${APP} main.go

.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/${APP}-mac main.go

.PHONY: test
test:
	go test ./... -coverprofile out.prof
	go tool cover -func ./out.prof | grep total
	rm ./out.prof

.PHONY: coverage
coverage:
	go test ./... -coverprofile out.prof
	go tool cover -html=out.prof
	rm ./out.prof
