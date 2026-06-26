# kaliki-template

A [Copier](https://copier.readthedocs.io/) template for the Go + Flutter +
OpenAPI stack, Docker-first.

## Generate a project

```sh
uvx copier copy gh:kaliki-tech-labs/kaliki-template my-app
# or from a local checkout:
uvx copier copy . ../my-app
```

You'll be prompted for identity (name, Go module path, bundle id, db name) and
module choices (auth, caching, eventing, payments, push, object storage,
secondary frontend, example domain). Observability is always included with a
no-op default.

## What you get

- **Backend** — Go (gin + pgx), config-from-YAML, auto-migrations on startup,
  first-class observability (no-op until configured), one example domain.
- **API** — OpenAPI specs as the contract source of truth.
- **Frontend** — Flutter (Riverpod + GoRouter); `shared_ui` vendored or as a
  git dependency.
- **Docker** — `docker compose up --build` brings up postgres (+redis/redpanda
  per flags) and the backend together.

## Test the template

```sh
./test-template.sh             # renders combos, go vet, docker smoke test
SKIP_DOCKER=1 ./test-template.sh   # render + go vet only (fast)
```

## Status

Basic structuring. Codegen pipeline (sqlc / oapi-codegen / redocly build /
Dart client regen), CI matrix, and the Tier-2 module internals are scaffolded
and documented but not yet fully fleshed out — see `TODO` below.

### TODO
- [ ] sqlc + oapi-codegen wiring (replace hand-written store/handler in `items`)
- [ ] redocly combined-spec build + Dart client regen script
- [ ] flesh out auth (JWT/invite/reset), eventing (outbox), payments, push, storage
- [ ] GitHub Actions matrix that generates + builds combos
- [ ] extract shared_ui to its own repo (git-dependency mode)
