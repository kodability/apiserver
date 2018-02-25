FROM golang:1.10

LABEL maintainer="kodability"
LABEL name="kodability tryout runner"

RUN mkdir -p /go/src/github.com/kodability/tryout-runner
WORKDIR /go/src/github.com/kodability/tryout-runner
COPY . .

CMD "tryout-runner"
