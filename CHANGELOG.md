# Changelog

Versions are git tags; `copier update` upgrades a generated project between them.

## Unreleased

### Added
- **Frontend example list.** When the example domain is included, the home
  screen renders the `items` list with a read/write round-trip (`itemsClient`).
- **Frontend & API build on generation.** Post-gen tasks now build the combined
  OpenAPI spec (Node) and, for the frontend, generate the Dart API client,
  scaffold platform folders (`flutter create`), and run `build_runner` — the
  generated project is ready to `flutter run`.
- **CORS.** The API reflects any `Origin` in development (`APP_ENV != production`)
  and uses a `CORS_ALLOWED_ORIGINS` allowlist in production.
- CI now exercises `basic` and `token_mail` combos and runs `flutter analyze`;
  `test-template.sh` renders the working tree so local edits are tested
  pre-commit.

### Changed (breaking for `copier update`)
- **OTP signup no longer starts a session** (`auth=jwt-full`, OTP method). Signup
  returns a short-lived `verification_token` (`SignupResponse`) instead of an
  access/refresh pair; the token is redeemed with the emailed code at
  `/auth/verify` to obtain the session. A new OTP screen implements the flow.
  Verification tokens cannot authenticate API calls.
- **`backend/config/env.yaml` is now generated as `env.example.yaml`** and copied
  to the git-ignored `env.yaml` on generation. After `copier update`, reconcile
  your live `env.yaml` against the new example.
- **New generation prerequisites:** Node/npm (API spec) and the Flutter SDK
  (frontend codegen). Previously only Go was required.

## v1.0.0

Initial release: Go (gin + pgx) backend, Flutter (Riverpod + GoRouter) frontend,
OpenAPI-driven codegen, auth (JWT signup/login, verification, password reset),
Tier-2 module integration points, observability, and a Docker-first stack.
