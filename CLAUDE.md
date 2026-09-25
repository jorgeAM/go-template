# go-template

Go 1.26 modular monolith template: hexagonal modules under `internal/<module>/`, shared I/O
adapters in `internal/platform/`, domain-agnostic types in `internal/shared/`. Chi HTTP server,
PostgreSQL via pgx + sqlc, OpenAPI via oapi-codegen, mocks via gomock.

Detailed conventions live in `.claude/rules/` — this file is the map, the rules are the law.

## Modules

| Module | Path | Responsibility | HTTP | Schema |
|--------|------|----------------|------|--------|
| identity | `internal/identity/` | Users: create and fetch user accounts (name, email, bcrypt-hashed password) | `POST /api/v1/user`, `GET /api/v1/user/{id}` | `identity` |

`internal/identity/` is the reference implementation every new module copies.

## Boot Flow

1. `internal/bootstrap.Modules()` returns the module list (dependency-ordered).
2. `cmd/app` calls each module's `Init(ctx)` (module builds its own DB pool, repositories,
   platform adapters from its own `config.go`), then `RegisterHttp(ctx, router)`.
3. `cmd/migrate` applies each module's embedded `MigrationFS()`, tracked per module in
   `schema_migrations_<module>`.

A new module implements `internal/shared/module.Module` (`Name`, `Init`, `RegisterHttp`,
`MigrationFS`) and is added to `bootstrap.Modules()` — nothing else is wired centrally.

## Module Layout

```
internal/<module>/
├── module.go        # module.Module implementation, owns all wiring
├── config.go        # env vars this module reads
├── domain/          # aggregates, invariants, sentinel errors, repository interfaces, unmarshall.go
├── app/             # use cases: command/, query/, models/ (output shapes)
├── api/http/        # openapi.yaml → oapi-codegen strict server, server.go
├── adapters/db/     # migrations/, queries/, sqlc.yaml → sqlc/, postgres_<name>_repository.go
└── mocks/           # gomock mocks of this module's own interfaces
```

See `.claude/rules/folder-structure.md` and `.claude/rules/domain-layer.md`.

## Shared vs Platform

- `internal/platform/` — I/O adapters: `db` (pgx transactor), `eventbus` (in-memory, SNS, SQS),
  `http` (middleware, health handler, REST client), `log` (zap-backed structured logging),
  `mailer` (SES, SendGrid, in-memory), `storage` (Cloudflare R2 presigned URLs).
  See `.claude/rules/platform-folder-structure.md`.
- `internal/shared/` — zero-dependency types: `errors` (sentinel `ErrorCode` + `Error`),
  `valueobject` (`ID`, `Email`, `Timestamps`), `criteria`, `events`, `module`, `crypto`, `env`,
  `collections`, `generator`, `ref`. See `.claude/rules/shared-folder-structure.md`.

## Commands

| Command | What it does |
|---------|--------------|
| `make run` | Start the server (`cmd/app`), logs piped through `jq` |
| `make generate` | `go generate ./...`: sqlc, oapi-codegen, mockgen (all via `go tool`) |
| `make unit-tests` | `go test ./...` with coverage |
| `make migrate` | Apply every module's pending migrations (`cmd/migrate`) |
| `make new_migration MODULE=<m> MIGRATION_NAME=<n>` | New migration pair for a module (golang-migrate CLI) |
| `make migration_down MODULE=<m>` | Roll back a module's last migration |
| `make tidy` | `go mod tidy` + `go mod vendor` |

Never hand-edit generated code (`sqlc/`, `openapi.gen.go`, `mocks/`); change the source and run
`make generate`.

## Configuration

Env vars (see `.env.example`): `APP_ENV`, `PORT`, `POSTGRES_HOST`, `POSTGRES_PORT`,
`POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_MAX_OPEN_CONNECTIONS`, `JWT_KEY`,
`JWT_ISSUER`. `.env` is loaded automatically and never committed.

## Rules Index

| Topic | Rule |
|-------|------|
| Where code goes | `.claude/rules/folder-structure.md` |
| Aggregates, repositories, unmarshall | `.claude/rules/domain-layer.md` |
| Errors, handlers, comments, Go style | `.claude/rules/golang.md` |
| Logging | `.claude/rules/logging.md` |
| Tests | `.claude/rules/testing.md` |
| Branches, commits, secrets | `.claude/rules/git-workflow.md` |
| `internal/platform/` admission | `.claude/rules/platform-folder-structure.md` |
| `internal/shared/` admission | `.claude/rules/shared-folder-structure.md` |
