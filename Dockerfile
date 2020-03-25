FROM golang:alpine as build-env

ENV GO111MODULE=on
RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /chatroom
WORKDIR /chatroom
COPY chatroom /chatroom

RUN go mod download
RUN go build -o chatroom
CMD ./chatroom server
