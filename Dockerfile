FROM golang:1.12-alpine AS builder

LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"
COPY thumbnailer.go /go/src/thumbnailer.go
WORKDIR /go/src
RUN go build -o /go/bin/thumbnailer /go/src/thumbnailer.go

FROM alpine:latest
RUN apk add imagemagick
COPY --from=builder /go/bin/thumbnailer /bin/thumbnailer
ENV PORT=1337
EXPOSE $PORT
ENTRYPOINT ["/bin/thumbnailer"]
