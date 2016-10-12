FROM golang:latest
#FROM resin/raspberrypi3-golang

# Copy the local package files to the container’s workspace.
WORKDIR /gopath/src/github.com/KanybekMomukeyev/goDockerCompose/api

ADD . /gopath/src/github.com/KanybekMomukeyev/goDockerCompose/

ADD views /gopath/src/github.com/KanybekMomukeyev/goDockerCompose/api/views
ADD conf /gopath/src/github.com/KanybekMomukeyev/goDockerCompose/api/conf
ADD controllers /gopath/src/github.com/KanybekMomukeyev/goDockerCompose/api/controllers
ADD models /gopath/src/github.com/KanybekMomukeyev/goDockerCompose/api/models

# Install our dependencies
RUN go get github.com/astaxie/beego
RUN go get github.com/beego/bee
RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx

RUN go get github.com/KanybekMomukeyev/goDockerCompose/api

RUN go install github.com/KanybekMomukeyev/goDockerCompose/api

#ADD api/server.go /gopath/src/github.com/KanybekMomukeyev/goDockerCompose/api/server.go

EXPOSE 8080

ENTRYPOINT /go/bin/api

#CMD ["go", "run", "/gopath/src/github.com/KanybekMomukeyev/goDockerCompose/api/server.go"]