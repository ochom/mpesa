# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.22-alpine AS builder

RUN apk add build-base

WORKDIR /app

COPY go.* ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -o ./tmp/server cmd/server/main.go

##
## Deploy
##
FROM busybox:1.35.0-uclibc AS deploy 

WORKDIR /

COPY --from=builder /app/tmp/server /server

RUN mkdir -p /data
RUN mkdir -p /data/certs

EXPOSE 8080

ENTRYPOINT [ "./server" ]
