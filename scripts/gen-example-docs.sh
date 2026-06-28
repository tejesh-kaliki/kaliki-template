#!/usr/bin/env bash
# Generate examples/README.md from examples/profiles/*.yml.
#
# Single source of truth = the profile files. This emits a coverage matrix and a
# per-profile section. CI runs this and fails on a non-empty git diff, so the
# docs can never drift from the actual profiles. Never edit examples/README.md
# by hand.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
exec python3 "$ROOT/scripts/_gen_example_docs.py" "$@"
