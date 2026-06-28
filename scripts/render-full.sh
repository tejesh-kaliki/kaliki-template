#!/usr/bin/env bash
# Render COMPLETE, runnable repos into examples/_full/<profile>/ for local
# browsing. Unlike the committed pre-task snapshots (examples/rendered/), these
# run copier _tasks (backend codegen, go mod tidy) AND the frontend bootstrap
# (API client + flutter create + build_runner), so each is a real project you can
# `cd` into and run. examples/_full/ is gitignored — it is scratch, never
# committed. Use validate-examples.sh if you also want build/test.
#
# Usage: scripts/render-full.sh [profile ...]   (default: all)
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROFILES_DIR="$ROOT/examples/profiles"
FULL_DIR="$ROOT/examples/_full"

has_frontend() { python3 -c "import yaml,sys; print(yaml.safe_load(open(sys.argv[1]))['answers'].get('include_frontend', False))" "$1"; }

profiles=("$@")
if [ ${#profiles[@]} -eq 0 ]; then
  for f in "$PROFILES_DIR"/*.yml; do profiles+=("$(basename "$f" .yml)"); done
fi

for name in "${profiles[@]}"; do
  src="$PROFILES_DIR/$name.yml"
  dst="$FULL_DIR/$name"
  echo ">> rendering full $name -> ${dst#$ROOT/}"

  data_file="$(mktemp)"
  python3 -c "import yaml,sys; yaml.safe_dump(yaml.safe_load(open(sys.argv[1]))['answers'], open(sys.argv[2],'w'), sort_keys=False)" "$src" "$data_file"

  rm -rf "$dst"
  mkdir -p "$FULL_DIR"
  # WITH tasks (--trust) — this is the post-task tree.
  copier copy --defaults --trust --data-file "$data_file" --vcs-ref HEAD "$ROOT" "$dst"

  if [ "$(has_frontend "$src")" = "True" ]; then
    ( cd "$dst" && bash scripts/regenerate-api-client.sh )
    pname="$(awk '/^name:/{print $2; exit}' "$dst/frontend/pubspec.yaml")"
    ( cd "$dst/frontend" \
        && flutter create --project-name "$pname" . \
        && rm -f test/widget_test.dart \
        && flutter pub get \
        && dart run build_runner build --delete-conflicting-outputs )
  fi
  rm -f "$data_file"
done

echo "done. browse examples/_full/<profile>/ (gitignored scratch)."
