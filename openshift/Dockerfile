FROM thumbnailer-base

LABEL maintainer="Patrick Bucher <patrick.bucher@stud.hslu.ch>"

LABEL io.openshift.s2i.destination="/tmp" \
      io.openshift.s2i.scripts-url="image:///home/gopher/s2i"

USER 1001

RUN mkdir /home/gopher/s2i
COPY ./s2i/bin /home/gopher/s2i

# already exposed?
# EXPOSE 1337

CMD ["/home/gopher/s2i/usage"]
