# Local validation targets. The CI free tier only runs a cheap drift check per PR
# (examples-drift.yml) and the heavy matrix nightly (validate.yml); the real
# build/test gate runs here on your machine. Run `make check` before pushing.

.PHONY: check snapshots docs validate hooks full

# Full local gate: refresh snapshots + docs, fail on drift, then build/test every
# profile's post-task render. This is what protects main between nightly runs.
check: snapshots docs
	@git diff --quiet || { echo "drift: commit regenerated snapshots/docs"; git --no-pager diff --stat; exit 1; }
	@$(MAKE) validate

# Pre-task snapshots (committed, browsable) + generated docs.
snapshots:
	bash scripts/render-examples.sh

docs:
	bash scripts/gen-example-docs.sh

# Post-task build + test for every profile (the heavy bit; needs Go/Node/Flutter).
validate:
	bash scripts/validate-examples.sh

# Render complete, runnable repos into examples/_full/<profile>/ (gitignored
# scratch) for local browsing — post-task, with frontend codegen. `make full p=jwt-basic`
# renders one profile; no arg renders all.
full:
	bash scripts/render-full.sh $(p)

# Install the pre-push hook so `make check` runs automatically before each push.
hooks:
	git config core.hooksPath .githooks
	@echo "pre-push hook enabled (core.hooksPath=.githooks)"
