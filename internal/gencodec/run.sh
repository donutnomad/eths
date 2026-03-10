#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
BIN="$ROOT_DIR/bin/gencodec"

if [ ! -f "$BIN" ] || [ "$(find "$SCRIPT_DIR" -name '*.go' -newer "$BIN" 2>/dev/null)" ]; then
  cd "$SCRIPT_DIR"
  go build -o "$BIN" .
fi

exec "$BIN" "$@"
