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

# enable cgo for sqlite to work
ENV CGO_ENABLED=1 

# build the server
RUN go build -o /server .

# build the seeder
RUN go build -o /seeder ./cmd/seeder/main.go

##
## Deploy
##
FROM alpine:3.20.0 AS deploy

COPY --from=build /server /bin/app
COPY --from=build /seeder /bin/seeder

RUN chmod +x /bin/app
RUN chmod +x /bin/seeder

RUN mkdir -p /data
RUN mkdir -p /data/certs

EXPOSE 8080

CMD [ "/bin/app" ]
