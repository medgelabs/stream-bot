FROM golang:alpine as builder
WORKDIR /usr/src/app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum /usr/src/app/
RUN go mod download

COPY . /usr/src/app

ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -a -o medgebot main.go

# Final, clean image
FROM alpine
WORKDIR /app
COPY --from=builder /usr/src/app/medgebot /app/
COPY --from=builder /usr/src/app/config.yaml /app/

# Admin / API
EXPOSE 8080

ENV CHANNEL="medgelabs"
CMD /app/medgebot -channel $CHANNEL -all
