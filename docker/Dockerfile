FROM golang:1.23.6-alpine as builder

ENV config_env=docker

WORKDIR /app

COPY ./ /app

RUN go mod download


FROM golang:1.23.6-alpine as runner

WORKDIR /app

COPY --from=builder ./app ./app

RUN go install github.com/githubnemo/CompileDaemon@latest

ENV config_env=docker

EXPOSE 3000

ENTRYPOINT CompileDaemon --build="go build ./cmd/api/main.go" --command=./main

