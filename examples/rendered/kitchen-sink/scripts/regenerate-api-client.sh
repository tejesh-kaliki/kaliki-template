#!/usr/bin/env bash
# Rebuild the combined OpenAPI spec, then regenerate the Dart API client.
# Run after any change under api/services/.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> Building combined OpenAPI spec"
cd "$ROOT_DIR/api"
npm install --silent
npm run build:local

echo "==> Regenerating Dart client"
cd "$ROOT_DIR/packages/shared_api_client"
rm -rf lib
mkdir -p lib
dart pub get
dart run swagger_parser

echo "==> Copying hand-written overrides"
cp -r lib_custom/. lib/ 2>/dev/null || true

echo "==> build_runner"
dart run build_runner build --delete-conflicting-outputs

echo "==> Done"
