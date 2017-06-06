FROM golang:latest

WORKDIR /go/src/github.com/KanybekMomukeyev/goDockerCompose/api

ADD . /go/src/github.com/KanybekMomukeyev/goDockerCompose/

ADD views /go/src/github.com/KanybekMomukeyev/goDockerCompose/api/views
ADD models /go/src/github.com/KanybekMomukeyev/goDockerCompose/api/models

ADD server.go /go/src/github.com/KanybekMomukeyev/goDockerCompose/api/server.go

# Install our dependencies
RUN go get github.com/kataras/iris
RUN go get github.com/Sirupsen/logrus
RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx

RUN go install github.com/KanybekMomukeyev/goDockerCompose/api

EXPOSE 8080

ENTRYPOINT /go/bin/api

#CMD ["go", "run", "/go/src/github.com/KanybekMomukeyev/goDockerCompose/api/server.go"]