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
module choices (auth, mailer, verification method, caching, eventing, payments,
push, object storage, example domain). Observability is always included with a
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
- [x] auth: JWT signup/login, email verification + password reset, argon2id,
      refresh-token rotation + revocation, and a protected `GET /auth/me`
- [x] transactional email (`mailer` flag): SMTP + SES + log transports, Mailpit
      in the local stack; `verification_method` flag (OTP code vs token link).
      Credentials are emailed/logged, never returned in API responses
- [x] Tier-2 modules: caching (redis), eventing (kafka + outbox), payments,
      push, object storage — integration points, gated by flags
- [x] transactional outbox relay worker (`cmd/workers`) for eventing
- [x] frontend auth flow: secure token storage, refresh interceptor, Riverpod
      auth controller, GoRouter redirect, login/signup/home screens, local
      widgets (no shared_ui package — extract one by hand if a 2nd app needs it)
- [x] API client generated via a swagger_parser fork supporting `x-dart-type`
      (`format: date` → `LocalDate`, not `DateTime`)
- [x] GitHub Actions matrix (renders + builds + tests minimal/default/full)
- [x] integration-style handler tests for `items` + `auth` (real DB, worked
      examples in `internal/testsupport`)

### Possible next
- OTP / auth-endpoint brute-force protection (rate limiting)
- a worked example of a domain route behind `auth.Middleware()` (the middleware
  and `/auth/me` exist; `items` is intentionally left public)
