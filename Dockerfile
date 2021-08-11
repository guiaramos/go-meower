# Build stage
FROM golang:1.10.2-alpine3.7 AS build
# Support CGO and SSL
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/guiaramos/meower

# Copy each service
COPY util util
COPY event event
COPY db db
COPY search search
COPY schema schema
COPY meow_service meow_service
COPY query_service query_service
COPY pusher_service pusher_service

# Compile them
RUN go install ./...

# Production build stage
FROM alpine:3.7
WORKDIR /usr/bin
# Copy built binaries
COPY --from=build /go/bin .
