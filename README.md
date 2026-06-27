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

Working stack. Codegen pipeline, auth, Tier-2 modules, and CI are in place and
verified end-to-end (see `test-template.sh` / `.github/workflows/template-ci.yml`).

### Done
- [x] sqlc + oapi-codegen wiring (`items` is a full generated vertical slice)
- [x] redocly auto-discovery spec build + Dart client regen script
- [x] auth: JWT signup/login, email verification + password reset, argon2id
- [x] transactional email (`mailer` flag): SMTP + SES + log transports, Mailpit
      in the local stack; `verification_method` flag (OTP code vs token link).
      Credentials are emailed/logged, never returned in API responses
- [x] Tier-2 modules: caching (redis), eventing (kafka + outbox), payments,
      push, object storage — integration points, gated by flags
- [x] GitHub Actions matrix (renders + builds + tests minimal/default/full)
- [x] integration-style handler tests for `items` + `auth` (real DB, worked
      examples in `internal/testsupport`)

### Not planned (for now)
- shared_ui repo extraction — `shared_ui_source=git` flag exists; defer the
  actual split until a second product consumes it.

### Possible next
- transactional outbox relay worker (cmd/workers) for eventing
- a worked example of a JWT-protected route (auth middleware exists but no
  endpoint demonstrates it yet)
- OTP brute-force protection (attempt counter / rate limiting)
- refresh-token / revocation flow (JWT is currently stateless HS256, no refresh)
- frontend auth screens (frontend is currently scaffold + generated client only)
