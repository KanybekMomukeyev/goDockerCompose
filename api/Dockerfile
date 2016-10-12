#FROM golang:latest
FROM resin/raspberrypi3-golang

# Copy the local package files to the container’s workspace.
WORKDIR /go/src/github.com/KanybekMomukeyev/goDockerCompose/api

ADD . /go/src/github.com/KanybekMomukeyev/goDockerCompose/

ADD views /go/src/github.com/KanybekMomukeyev/goDockerCompose/api/views
ADD conf /go/src/github.com/KanybekMomukeyev/goDockerCompose/api/conf
ADD controllers /go/src/github.com/KanybekMomukeyev/goDockerCompose/api/controllers
ADD models /go/src/github.com/KanybekMomukeyev/goDockerCompose/api/models
ADD server.go /go/src/github.com/KanybekMomukeyev/goDockerCompose/api/server.go

# Install our dependencies
RUN go get github.com/astaxie/beego
RUN go get github.com/beego/bee
RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx

RUN go install github.com/KanybekMomukeyev/goDockerCompose/api

EXPOSE 8080

ENTRYPOINT /go/bin/api

#CMD ["go", "run", "/go/src/github.com/KanybekMomukeyev/goDockerCompose/api/server.go"]