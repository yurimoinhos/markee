#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NEXT_DIR="${ROOT_DIR}/frontend-next"
NEXT_PORT="${NEXT_PORT:-3001}"
BACKEND_ADDR="${SERVER_ADDR:-:8000}"

export NEXT_INTERNAL_URL="${NEXT_INTERNAL_URL:-http://127.0.0.1:${NEXT_PORT}}"

cleanup() {
  if [[ -n "${NEXT_PID:-}" ]]; then
    kill "${NEXT_PID}" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT INT TERM

echo "[INFO] Iniciando Next.js em :${NEXT_PORT}"
(
  cd "${NEXT_DIR}"
  npm run dev -- --port "${NEXT_PORT}"
) &
NEXT_PID=$!

sleep 2

echo "[INFO] Iniciando backend Go em ${BACKEND_ADDR} com NEXT_INTERNAL_URL=${NEXT_INTERNAL_URL}"
cd "${ROOT_DIR}"
go run .
