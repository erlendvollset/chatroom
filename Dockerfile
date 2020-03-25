FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /chatroom
RUN mkdir -p /chatroom/proto

WORKDIR /chatroom

COPY ./proto/chatroom.pb.go /chatroom/proto
COPY ./main.go /chatroom

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o chatroom

CMD ./chatroom
