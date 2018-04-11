FROM golang:1.9.5-alpine
MAINTAINER muslimi

ENV SOURCES /go/src/github.com/u/muslimi/go-cloud/

COPY . ${SOURCES}

RUN cd ${SOURCES} && CGO_ENABLED=0 go install

ENV PORT 6060
EXPOSE 6060

ENTRYPOINT go-cloud
