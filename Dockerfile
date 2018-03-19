FROM golang:1.10

LABEL maintainer="kodability"
LABEL name="kodability API server"

RUN mkdir -p /go/src/github.com/kodability/apiserver
WORKDIR /go/src/github.com/kodability/apiserver
COPY . .

CMD "apiserver"
