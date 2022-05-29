FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine

COPY --from=builder /build/main /
COPY --from=builder /build/configs/config.json /configs/config.json
COPY --from=builder /build/web /web
COPY --from=builder /build/scripts /scripts

EXPOSE 8080

ENTRYPOINT ["./main"]