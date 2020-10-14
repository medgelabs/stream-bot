FROM golang:alpine as builder
WORKDIR /usr/src/app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum /usr/src/app/
RUN go mod download

COPY . /usr/src/app

# COPY main.go /usr/src/app/main.go
# COPY ./irc /usr/src/app/irc
# COPY ./bot /usr/src/app/bot
# COPY ./secret /usr/src/app/secret

ENV GOOS=linux
ENV CGO_ENABLED=0
RUN go build -a -o medgebot main.go

# Final, clean image
FROM alpine
WORKDIR /app
COPY --from=builder /usr/src/app/medgebot /app/
CMD ["/app/medgebot"]
