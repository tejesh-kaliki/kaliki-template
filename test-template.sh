#!/usr/bin/env bash
# Local mirror of template-ci: render representative combos, build + test the
# backend, and build the API spec for each. Docker smoke-tests the default combo.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

# name|flags
COMBOS=(
  "minimal|-d auth=none -d include_example_domain=false -d include_frontend=false -d caching=none"
  "default|"
  "full|-d eventing=kafka-redpanda -d payments=razorpay -d push_notifications=firebase -d object_storage=s3 -d include_secondary_frontend=true"
)

for combo in "${COMBOS[@]}"; do
  name="${combo%%|*}"; flags="${combo#*|}"
  echo "==> render: $name"
  uvx copier copy --defaults --trust $flags "$ROOT" "$WORK/$name"

  echo "==> backend build + test: $name"
  ( cd "$WORK/$name/backend" && go build ./... && go test ./... )

  echo "==> api spec build: $name"
  ( cd "$WORK/$name/api" && npm install --silent && npm run build:local )
done

if [[ "${SKIP_DOCKER:-}" != "1" ]]; then
  echo "==> docker compose up (default) + smoke endpoints"
  cd "$WORK/default"
  cp .env.example .env
  docker compose up --build -d
  for i in $(seq 1 40); do curl -fsS localhost:8080/health >/dev/null 2>&1 && break; sleep 2; done
  echo "-- health:"; curl -fsS localhost:8080/health; echo
  docker compose down -v
else
  echo "==> skipped docker (SKIP_DOCKER=1)"
fi

echo "==> OK"
