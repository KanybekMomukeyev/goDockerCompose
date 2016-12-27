FROM golang:latest

# Copy the local package files to the container’s workspace.
WORKDIR /go/src/github.com/KanybekMomukeyev/goDockerCompose/grpc

ADD . /go/src/github.com/KanybekMomukeyev/goDockerCompose/

ADD client /go/src/github.com/KanybekMomukeyev/goDockerCompose/grpc/client
ADD model /go/src/github.com/KanybekMomukeyev/goDockerCompose/grpc/model
ADD proto /go/src/github.com/KanybekMomukeyev/goDockerCompose/grpc/proto
ADD server /go/src/github.com/KanybekMomukeyev/goDockerCompose/grpc/server
ADD mainserver.go /go/src/github.com/KanybekMomukeyev/goDockerCompose/grpc/mainserver.go

# Install our dependencies
RUN go get github.com/astaxie/beego
RUN go get github.com/beego/bee
RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx
RUN go get golang.org/x/net/context
RUN go get github.com/golang/protobuf
RUN go get google.golang.org/grpc



RUN go install github.com/KanybekMomukeyev/goDockerCompose/grpc

EXPOSE 8080

ENTRYPOINT /go/bin/api

#CMD ["go", "run", "/go/src/github.com/KanybekMomukeyev/goDockerCompose/api/server.go"]