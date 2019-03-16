FROM golang:alpine

LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"

LABEL io.openshift.s2i.destination="/tmp" \
      io.openshift.s2i.scripts-url="image:///home/gopher/s2i"

RUN apk add imagemagick && \
    addgroup --gid 1001 gophers && \
    adduser -D --uid 1001 -G gophers gopher

USER 1001

RUN mkdir /home/gopher/s2i && \
    mkdir /home/gopher/src && \
    mkdir /home/gopher/bin

COPY ./s2i/bin /home/gopher/s2i

EXPOSE 1337

CMD ["/home/gopher/s2i/usage"]
