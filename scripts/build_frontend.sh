#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="${ROOT_DIR}/frontend-next"
WEBUI_DIST="${ROOT_DIR}/webui/dist"

mkdir -p "${WEBUI_DIST}"

if ! command -v npm >/dev/null 2>&1; then
  echo "[ERROR] npm não encontrado no ambiente"
  exit 1
fi

pushd "${FRONTEND_DIR}" >/dev/null

if [[ ! -d node_modules ]]; then
  npm install
fi

npm run build

popd >/dev/null

cp "${FRONTEND_DIR}/public/next-unavailable.html" "${WEBUI_DIST}/index.html"

echo "[INFO] Frontend Next.js buildado com sucesso."
echo "[INFO] Fallback embed atualizado em ${WEBUI_DIST}/index.html"
