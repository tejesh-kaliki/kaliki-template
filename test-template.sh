#!/usr/bin/env bash
# Smoke-test the template: render a few flag combinations, then build/run the
# default one end-to-end via Docker and hit the health + items endpoints.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

copier() { uvx copier "$@"; }

echo "==> [1/4] render: default"
copier copy --defaults --trust "$ROOT" "$WORK/default"

echo "==> [2/4] render: minimal (no auth, no example, no frontend)"
copier copy --defaults --trust \
  -d auth=none -d include_example_domain=false -d include_frontend=false \
  "$ROOT" "$WORK/minimal"

echo "==> [3/4] render: everything (secondary frontend, kafka, redis)"
copier copy --defaults --trust \
  -d include_secondary_frontend=true -d eventing=kafka-redpanda \
  "$ROOT" "$WORK/full"

echo "==> rendered trees:"
( cd "$WORK/default" && find . -path ./.git -prune -o -type f -print | sort )

echo "==> backend go vet (default)"
( cd "$WORK/default/backend" && go mod tidy && go vet ./... )

if [[ "${SKIP_DOCKER:-}" != "1" ]]; then
  echo "==> [4/4] docker compose up (default) + smoke endpoints"
  cd "$WORK/default"
  cp .env.example .env
  docker compose up --build -d
  # wait for health
  for i in $(seq 1 30); do
    if curl -fsS localhost:8080/health >/dev/null 2>&1; then break; fi
    sleep 2
  done
  echo "-- health:";  curl -fsS localhost:8080/health; echo
  echo "-- create:";  curl -fsS -X POST localhost:8080/api/v1/items \
      -H 'content-type: application/json' -d '{"name":"hello"}'; echo
  echo "-- list:";    curl -fsS localhost:8080/api/v1/items; echo
  docker compose down -v
else
  echo "==> [4/4] skipped docker (SKIP_DOCKER=1)"
fi

echo "==> OK"
