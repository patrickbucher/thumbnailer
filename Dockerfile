FROM golang:alpine

LABEL maintainer="patrick.bucher@stud.hslu.ch"
RUN apk add imagemagick

ENV SOURCE_DIR=/go/src/github.com/patrickbucher/thumbnailer
RUN mkdir -p "${SOURCE_DIR}"
COPY thumbnailer.go "${SOURCE_DIR}/"
WORKDIR ${SOURCE_DIR}

RUN go build -o /go/bin/thumbnailer thumbnailer.go

ENTRYPOINT /go/bin/thumbnailer

EXPOSE 1337
