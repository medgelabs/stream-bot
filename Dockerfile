FROM golang:buster
WORKDIR /usr/src/app

COPY go.mod go.sum /usr/src/app/

COPY main.go /usr/src/app/main.go
COPY ./irc /usr/src/app/irc
COPY ./bot /usr/src/app/bot
ENV GOOS=linux
RUN go build -o medgebot main.go

CMD ["/usr/src/app/medgebot"]
