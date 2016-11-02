#FROM golang:1.7-alpine
FROM resin/raspberrypi3-golang

WORKDIR /go/src/app

COPY . /go/src/app

RUN \
	#apk add --no-cache git && \
	apt-get install git-core && \
	go-wrapper download && \
	go-wrapper install -ldflags "$(go run buildscripts/gen-ldflags.go)" && \
	mkdir -p /export/docker && \
	cp /go/src/app/docs/Docker.md /export/docker/ && \
	rm -rf /go/pkg /go/src && \
	apk del git

EXPOSE 9000
ENTRYPOINT ["go-wrapper", "run"]
CMD ["server"]
VOLUME ["/export"]

#####

FROM resin/raspberrypi3-golang
MAINTAINER kano.momuke@gmail.com
ADD kanoexport /bin/
COPY minio /bin/
CMD ["/bin/minio", "server", "/kanoexport"]
#dockerfile from exe

######################################################
FROM resin/raspberrypi3-golang
MAINTAINER Kanybek Momukeev <kano.momuke@gmail.com>

ADD kanoexport /bin/
COPY minio /bin/

ENV PORT 9000
EXPOSE 9000

ENTRYPOINT ["/contributors"]
CMD ["/bin/minio", "server", "/kanoexport"]

######################################################
FROM google/debian:wheezy
MAINTAINER Kelsey Hightower <kelsey.hightower@gmail.com>
ADD contributors contributors
ENV PORT 80
EXPOSE 80
ENTRYPOINT ["/contributors"]
######################################################

FROM google/debian:wheezy
MAINTAINER Kanybek Momukeev <kano.momuke@gmail.com>
ENV EXPORT_PATH "/kanoexport"
ADD minio minio
ADD kanoexport kanoexport
EXPOSE 9000
#ENTRYPOINT ["/minio server"]
CMD ["/minio", "server", "/kanoexport"]