FROM golang:1.26.4-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY main.go ./
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/bar-lobby-protocol-service .

FROM alpine:3.22

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=build /out/bar-lobby-protocol-service /app/bar-lobby-protocol-service
COPY assets ./assets

USER app

ENV ADDR=:47777
EXPOSE 47777

ENTRYPOINT ["/app/bar-lobby-protocol-service"]
