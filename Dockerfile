FROM golang:latest
RUN mkdir -p /go/src/github.com/dcoxall
ADD . /go/src/github.com/dcoxall/juggler
WORKDIR /go/src/github.com/dcoxall/juggler
ENV GOBIN=/usr/local/bin
