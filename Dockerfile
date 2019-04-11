FROM golang:1.12.3-stretch AS builder

LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"

RUN apt-get update && apt-get install -y git ca-certificates
COPY thumbnailer.go go.mod /src/
WORKDIR /src
RUN go build -o /app/thumbnailer thumbnailer.go

FROM debian:stretch-slim
RUN apt-get update && apt-get install -y imagemagick
COPY --from=builder /app/thumbnailer /bin/thumbnailer
ENV PORT=1337
EXPOSE $PORT
CMD ["/bin/thumbnailer"]
