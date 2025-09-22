# ----- Stage 1: Build the Go application -----
ARG TARGETPLATFORM
ARG BUILDPLATFORM
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder
WORKDIR /app

# 依存キャッシュを有効化
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# ソース投入
COPY . .

# buildx が注入するターゲットOS/ARCHを利用（マルチアーキ対応）
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" -a -installsuffix cgo \
    -o /app/main ./cmd/main.go

# ----- Stage 1.5: Bring migrate binary (arch-aware) -----
FROM --platform=$TARGETPLATFORM migrate/migrate:v4.17.0 AS migrator

# ----- Stage 2: Final runtime image -----
FROM --platform=$TARGETPLATFORM alpine:3.20
WORKDIR /app

RUN apk add --no-cache tzdata bash curl postgresql-client

# アプリ本体
COPY --from=builder /app/main /app/main

# golang-migrate のバイナリを取り込み
COPY --from=migrator /usr/local/bin/migrate /usr/local/bin/migrate

# migrationファイルをコピー
COPY db/migrations /app/migrations

# エントリポイント
COPY ./deploy/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENV TZ=Asia/Tokyo
EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
