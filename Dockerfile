FROM golang:buster as builder
WORKDIR /usr/src/app

COPY go.mod go.sum /usr/src/app/
COPY main.go /usr/src/app/main.go
COPY ./irc /usr/src/app/irc
COPY ./bot /usr/src/app/bot

ENV GOOS=linux
RUN go build -o medgebot main.go

# Final, clean image
FROM alpine
COPY --from=builder /usr/src/app/medgebot /app/medgebot
CMD ["/app/medgebot"]
