## Build
FROM golang:1.19-bullseye AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN mkdir build && go build -o build/appname cmd/server/main.go

## Deploy
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates

WORKDIR /app

COPY --from=build /app/build /app
COPY internal/infra/database/migrations /app/migrations
RUN chmod +x /app/appname && chown -R nobody:nogroup .

EXPOSE 8080

USER nobody:nogroup

ENTRYPOINT ["/app/appname"]