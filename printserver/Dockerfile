FROM golang:latest

WORKDIR /go/src/github.com/KanybekMomukeyev/goDockerCompose/printserver

ADD . /go/src/github.com/KanybekMomukeyev/goDockerCompose/

ADD build /go/src/github.com/KanybekMomukeyev/goDockerCompose/printserver/build
ADD server.go /go/src/github.com/KanybekMomukeyev/goDockerCompose/printserver/server.go

RUN go install github.com/KanybekMomukeyev/goDockerCompose/printserver

EXPOSE 8100

ENTRYPOINT /go/bin/printserver
