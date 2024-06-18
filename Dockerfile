# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.22-alpine  AS build

WORKDIR /app

COPY go.* ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN apk add --no-cache gcc g++ git openssh-client
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /server cmd/server/main.go



##
## Deploy
##
FROM busybox:1.35.0-uclibc AS deploy 

WORKDIR /

COPY --from=build /server .

RUN mkdir -p /data
RUN mkdir -p /data/certs

EXPOSE 8080

CMD ["/server"]
