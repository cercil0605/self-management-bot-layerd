# ----- Stage 1: Build the Go application -----
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 静的リンク（Alpineで動かす前提）
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o /app/main ./cmd/main.go

# ----- Stage 1.5: Bring migrate binary -----
FROM migrate/migrate:v4.17.0 AS migrator

# ----- Stage 2: Final runtime image -----
FROM alpine:latest
WORKDIR /app

# ルートCAが必要（DB/TLSやGHCRアクセス時に便利）
RUN apk add --no-cache ca-certificates tzdata bash curl postgresql-client

# アプリ本体
COPY --from=builder /app/main /app/main

# golang-migrate のバイナリを取り込み
COPY --from=migrator /usr/local/bin/migrate /usr/local/bin/migrate

COPY db/migrations /app/migrations

COPY ./deploy/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENV TZ=Asia/Tokyo
EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
