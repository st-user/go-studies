# Multi stage building strategy for reducing image size.
FROM golang:1.19.3-alpine3.16 AS build-env
# ENV GO111MODULE=on
RUN mkdir /app
WORKDIR /app

# Install each dependencies
COPY go.mod /app
# COPY go.sum /app
RUN go mod download
RUN apk add --no-cache --virtual .builddeps git gcc make build-base alpine-sdk

# COPY main module
COPY . /app

# Check and Build
RUN go get golang.org/x/lint/golint && \
    make validate && \
    make build-linux && \
	apk del .builddeps

### If use TLS connection in container, add ca-certificates following command.
### > RUN apk add --no-cache ca-certificates
FROM alpine:3.16.3
COPY --from=build-env /app/main /
EXPOSE 80
ENTRYPOINT ["/main"]
