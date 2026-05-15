#!/usr/bin/env bash

set -euo pipefail

OUTPUT_DIR="./dist"
PACKAGE="cmd/client/main.go"
LDFLAGS="-s -w"

mkdir -p "$OUTPUT_DIR"

echo "Building hours client..."
go build -ldflags="$LDFLAGS" -o "${OUTPUT_DIR}/hours" "$PACKAGE"
echo "Done. Binary written to ${OUTPUT_DIR}/hours"
