INSERT INTO public.dockerfile_templates
("id", "key", "name", "content")
VALUES(1, 'go_web', 'Go Web-server (with go.mod)', 'FROM golang:latest AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM scratch

COPY --from=builder /build/main /

EXPOSE 8080

ENTRYPOINT ["/main"]') ON CONFLICT DO NOTHING;

INSERT INTO public.dockerfile_templates
("id", "key", "name", "content")
VALUES(2, 'go_console', 'Go console (without go.mod)', 'FROM golang:latest AS builder

WORKDIR /build

COPY . .

RUN GO111MODULE=off CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine

COPY --from=builder /build/main /

CMD ["sleep", "infinity"]') ON CONFLICT DO NOTHING;

INSERT INTO public.dockerfile_templates
("id", "key", "name", "content")
VALUES(3, 'js_console', 'JS console', 'FROM node:12

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY . .

CMD [ "node", "index.js" ]') ON CONFLICT DO NOTHING;