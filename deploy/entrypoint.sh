#!/usr/bin/env sh
set -e

# ===== 構成 =====
# どちらか一方でOK：
# A) DATABASE_URL を .env で直接渡す
# B) 個別のENV（DB_HOST/PORT/USER/PASSWORD/NAME）から組み立てる
if [ -z "$DATABASE_URL" ]; then
  : "${DB_HOST:=db}"
  : "${DB_PORT:=5432}"
  : "${DB_USER:=user}"
  : "${DB_PASSWORD:=password}"
  : "${DB_NAME:=self_management}"
  DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
fi

# ===== DB起動待ち（pg_isready）=====
echo "⏳ Waiting for database at ${DB_HOST}:${DB_PORT} ..."
until pg_isready -h "${DB_HOST:-db}" -p "${DB_PORT:-5432}" -U "${DB_USER:-user}" >/dev/null 2>&1; do
  sleep 1
done
echo "✅ DB is ready."

# ===== マイグレーション適用（idempotent）=====
echo "🚀 Running migrations..."
/usr/local/bin/migrate -path /app/migrations -database "${DATABASE_URL}" up

# ===== アプリ起動 =====
echo "🏁 Starting app..."
exec /app/main
