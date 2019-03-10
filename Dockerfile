FROM golang:alpine

LABEL maintainer="patrick.bucher@stud.hslu.ch"

RUN apk add imagemagick && \
    addgroup --gid 1001 gophers && \
    adduser -D --uid 1001 -G gophers gopher

USER 1001:1001

ENV SRC_DIR=/home/gopher/src
ENV BIN_DIR=/home/gopher/bin
RUN mkdir $SRC_DIR && mkdir $BIN_DIR
COPY thumbnailer.go $SRC_DIR/
WORKDIR $SRC_DIR
RUN go build -o $BIN_DIR/thumbnailer thumbnailer.go

ENTRYPOINT $BIN_DIR/thumbnailer

EXPOSE 1337
