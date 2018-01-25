FROM golang:1.9

LABEL maintainer="kodability"
LABEL name="kodability tryout runner"

RUN mkdir -p /go/src/tryout-runner
WORKDIR /go/src/tryout-runner
COPY . .

RUN go-wrapper download
RUN go-wrapper install

CMD "tryout-runner"
