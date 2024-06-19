# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.22-alpine AS build

RUN apk add build-base

WORKDIR /app

COPY go.* ./
COPY go.sum ./

RUN go mod download

COPY . ./

ENV CGO_ENABLED=1 

RUN go build -o /server .

##
## Deploy
##
FROM busybox:1.35.0-uclibc AS deploy 

COPY --from=build /server /bin/app

RUN mkdir -p /data
RUN mkdir -p /data/certs

EXPOSE 8080

CMD [ "/bin/app" ]
