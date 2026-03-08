#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MOBILE_DIR="${ROOT_DIR}/mobile-native"

cd "${MOBILE_DIR}"
if [[ ! -d node_modules ]]; then
  npm install
fi
npm run start
