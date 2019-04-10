FROM golang:1.12-alpine AS builder

LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"

RUN apk add --no-cache git ca-certificates
COPY cmd/* /src/
WORKDIR /src
RUN go build -o /app/thumbnailer thumbnailer.go

FROM alpine:latest
RUN apk add imagemagick
COPY --from=builder /app/thumbnailer /bin/thumbnailer
ENV PORT=1337
EXPOSE $PORT
CMD ["/bin/thumbnailer"]
