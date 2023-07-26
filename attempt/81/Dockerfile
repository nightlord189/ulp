FROM golang:latest AS builder

WORKDIR /build

COPY . .

RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine

COPY --from=builder /build/main /