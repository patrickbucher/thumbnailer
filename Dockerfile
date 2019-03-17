FROM golang:alpine

LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"

RUN apk add imagemagick
COPY thumbnailer.go /go/src/
RUN go build -o /go/bin/thumbnailer /go/src/thumbnailer.go

EXPOSE 1337

ENTRYPOINT ["/go/bin/thumbnailer"]
