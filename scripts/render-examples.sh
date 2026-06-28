#!/usr/bin/env bash
# Render the committed, browsable snapshots in examples/rendered/<profile>/.
#
# These are PRE-TASK expansions (copier with _tasks skipped): pure template
# output, deterministic, the documentation surface. Volatile / task-generated
# files are stripped so the snapshot stays small and diffs stay meaningful.
# Full build + test validation happens separately in validate-examples.sh.
#
# Usage: scripts/render-examples.sh [profile ...]   (default: all)
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROFILES_DIR="$ROOT/examples/profiles"
RENDERED_DIR="$ROOT/examples/rendered"

# Files/dirs that are task outputs or otherwise volatile. They must never land in
# a committed snapshot; defense-in-depth alongside examples/rendered/.gitignore.
STRIP_GLOBS=(
  # Non-deterministic: _commit changes every commit, _src_path is the machine's
  # absolute path. Useless in a doc snapshot and would make the drift gate flap.
  ".copier-answers.yml"
  "backend/go.sum"
  "backend/config/env.yaml"
  "backend/gen"
  "frontend/pubspec.lock"
  "frontend/android" "frontend/ios" "frontend/linux"
  "frontend/macos" "frontend/windows" "frontend/web"
  "frontend/**/*.g.dart"
  "api/node_modules"
)

profiles=("$@")
if [ ${#profiles[@]} -eq 0 ]; then
  for f in "$PROFILES_DIR"/*.yml; do profiles+=("$(basename "$f" .yml)"); done
fi

extract_answers() {  # profile.yml -> a temp copier --data-file
  python3 - "$1" <<'PY'
import sys, yaml
data = yaml.safe_load(open(sys.argv[1]))
yaml.safe_dump(data["answers"], sys.stdout, sort_keys=False)
PY
}

for name in "${profiles[@]}"; do
  src="$PROFILES_DIR/$name.yml"
  dst="$RENDERED_DIR/$name"
  echo ">> rendering $name"
  data_file="$(mktemp)"
  extract_answers "$src" > "$data_file"

  rm -rf "$dst"
  # --defaults: don't prompt; --data-file: the profile answers.
  # -d _skip_tasks=true gates every _task off (see copier.yml) so this stays a
  # pure template expansion — deterministic, no Go/Node/Flutter needed.
  copier copy --defaults --trust -d _skip_tasks=true --data-file "$data_file" \
    --vcs-ref HEAD "$ROOT" "$dst"

  for g in "${STRIP_GLOBS[@]}"; do
    rm -rf "$dst"/$g 2>/dev/null || true
  done
  rm -f "$data_file"
done

echo "done. regenerate docs with scripts/gen-example-docs.sh"
