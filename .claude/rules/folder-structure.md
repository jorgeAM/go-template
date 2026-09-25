# Folder Structure

Rules for where code lives inside `internal/<module>/`. See root `CLAUDE.md` for module list and business responsibilities. Reference implementation: `internal/identity/`.

## Module Layout

Each business module follows this pattern:

```
internal/module/
├── module.go            # Bootstrap: implements shared module.Module (Name, Init, RegisterHttp).
│                         # Owns its lifecycle — constructs every platform concrete type this
│                         # module needs (DB pool, publisher, mailer, signer, ...) in Init(),
│                         # injects them as interfaces into app/ use cases.
├── config.go             # Module-owned config loader (env vars this module needs)
├── domain/                # Domain entities, value objects, business rules, repository interfaces
├── app/                    # Use cases and business logic
│   ├── command/              # Write operations
│   ├── query/                  # Read operations
│   └── models/                   # App-layer output shapes shared across command/query
├── api/                        # Driving (inbound) adapters — things that call INTO the module
│   └── http/                      # HTTP route handlers, OpenAPI-generated server/types
├── adapters/                        # Driven (outbound) adapters — things the module calls OUT to
│   └── db/                            # Repository implementations (sqlc, postgres, etc.)
├── client/                              # Interface + implementation for cross-module (in-process)
│                                         # communication — see "Cross-Module Client" below.
│                                         # ONLY when this module actually has a consumer today —
│                                         # do not scaffold an empty client/ "just in case".
└── mocks/                                  # Generated mocks for THIS module's OWN interfaces
                                             # (domain repositories). Never a mock of this module's
                                             # own client/ interface — see "Mock Ownership" below.
```

## Driving vs Driven

`api/` and `adapters/` are not interchangeable buckets — they answer different questions:

- `api/` — something calls **into** the module (HTTP request, ...). Primary/driving adapters.
- `adapters/` — the module calls **out** to something external (DB, ...). Secondary/driven adapters.

Keeping this split, instead of a single `infrastructure/` folder for both, preserves the actual
inbound/outbound distinction hexagonal architecture is built on.

## Standard Scaffolds

Two subfolders follow a fixed, codegen-driven scaffold — the same filenames and roles in every
module, only table/query/spec content differs. Copy the shape from `internal/identity/`, don't
reinvent it per module.

### `adapters/db/` (sqlc)

- `migrations/NNNNNN_<name>.up.sql` — this module's own migrations. Module-owned, not a shared
  migrations directory.
- `migrations.go` — `//go:embed migrations/*.sql` + `Migrations() fs.FS`. The module's
  `module.go` returns this from its `MigrationFS()` method (part of `module.Module`), so
  `make migrate` (`cmd/migrate`) picks the module up. Copy verbatim from any existing
  module. A module with no schema of its own (e.g. `notification`) has no `migrations.go`
  and returns `nil` from `MigrationFS()`.
- `queries/<name>.sql` — hand-written raw SQL, one file per aggregate.
- `sqlc.yaml` — sqlc config: `schema: migrations`, `queries: queries`, output package `sqlc`.
- `sqlc.go` — `//go:generate go tool sqlc generate` directive, nothing else in the file.
- `sqlc/` — fully generated (models, querier, query methods). Never hand-edit.
- `postgres_<name>_repository.go` — hand-written, satisfies `domain.<Name>Repository`, built on
  the generated `sqlc.Querier`.

### `api/http/` (oapi-codegen)

- `openapi.yaml` — hand-written OpenAPI spec, source of truth for the module's HTTP contract.
- `oapi-codegen.yaml` — codegen config (`chi-server: true`, `strict-server: true`, `models: true`).
- `openapi.go` — `//go:generate go tool oapi-codegen -config oapi-codegen.yaml openapi.yaml`
  directive, nothing else in the file.
- `openapi.gen.go` — fully generated: chi bindings, `StrictServerInterface`, typed
  `OpNNNJSONResponse` types. Never hand-edit.
- `server.go` — hand-written `Server` struct implementing `StrictServerInterface`, built from
  `app/` use cases in `NewServer(...)`, one handler method per operation (see golang.md's API
  Handlers section for the response pattern).

Both scaffolds regenerate via `make generate` (`go generate ./...`), same as mocks.

## Cross-Module Client

A module that exposes data or behavior to other modules does so through `client/` — not `api/`,
not `adapters/`:

- `client/client.go` — the interface (e.g. `IdentityClient`) plus its output DTOs (e.g.
  `AccountResponse`). Define DTO shapes here; never leak a `domain.*` type across the module
  boundary.
- `client/service.go` — the implementation, same package, calling into this module's own `app/`
  (or a `domain/` repository) in-process. Swappable later for a real HTTP/gRPC call without
  changing the interface or any consumer's import.

`client/` is not infrastructure: it performs no I/O — in-process call, no network, no
serialization — so it does not belong nested under `api/` or `adapters/` even though it is an
inbound-shaped concern from the consumer's point of view. It sits at module root as its own
concern.

Do not create `client/` speculatively. If no other module consumes this module today, there is
no `client/` folder — not even an empty one. Add it (interface + implementation together, per
above) the moment a real consumer shows up, not before.

### Mock Ownership

Mocks live wherever the tests that consume them live. This is one rule applied consistently:

- Mocks for this module's own `domain/*_repository.go` interfaces are consumed by this module's
  own `app/` tests → they live in this module's own `mocks/`.
- Mocks for this module's `client/client.go` interface are consumed by OTHER modules' tests → the
  producing module does **not** generate a mock for it. Each consumer module generates and owns
  its own mock, in its own `mocks/`, via a `//go:generate` directive declared in the consumer that
  points `-source` at the producer's `client/client.go`.

### Anemic Entities (No `domain/`)

Per [domain-layer.md](domain-layer.md) rule 8, an entity with no invariant enforcement (its
constructor validates nothing a raw struct literal couldn't, and it has no mutation methods)
moves down to `app/` instead of `domain/`. When that happens:

- The struct itself → `app/models/<name>.go`, same as any other output shape.
- Its repository interface → `app/models/<name>_repository.go`, colocated with the entity
  struct it serves — same reasoning as domain-layer.md rule 4 (repository lives next to what
  it repositories), just relocated to `models/` since that's where the anemic entity itself
  lives instead of `domain/`.
- Mock → still `mocks/`, `//go:generate` directive points `-source` at
  `app/models/<name>_repository.go`.
- `adapters/db/postgres_<name>_repository.go` is unaffected beyond its import path.

Hypothetical example: a `settlement` module's `app/models/order.go` + `order_repository.go` — a log record of
a paid order attached to a billing cycle, no invariant enforcement of its own beyond referencing
an existing cycle and order. Don't apply this preemptively; it only kicks in once an entity is
confirmed anemic per domain-layer.md rule 8.

## Service Integration Steps

1. **Interface Definition** — repository interfaces used only within the module go in
   `domain/[name]_repository.go`. The interface exposed to other modules goes in
   `client/client.go`.
2. **Implementation** — repository implementation goes in `adapters/db/[name]_repository.go`. The
   client implementation goes alongside its interface in `client/service.go`, calling this
   module's own `app/` in-process.
3. **Dependency Injection** — this module's own `module.go` `Init()` constructs every dependency
   this module owns (platform concretes, repositories, its `client.Service`) and exposes what
   other modules need as fields on `Module`. A consumer module wires the producer's `client`
   interface into its own `app/` layer.
4. **Mock Generation** — `//go:generate` comments, run `make generate` (this runs
   `go generate ./...` project-wide, so it doesn't matter which module's folder a directive lives
   in — see Mock Ownership above).

## Where New Code Goes

When creating or modifying modules:

1. **Domain First** — entities and business rules in `domain/`
2. **Application Logic** — use cases in `app/` (split into `command/` and `query/`)
3. **Inbound Adapters** — `api/` (e.g. `api/http/`)
4. **Outbound Adapters** — `adapters/` (e.g. `adapters/db/`)
5. **Cross-Module Exposure** — producer owns interface + implementation together in its own
   `client/`; consumer depends on the interface only, DI injects the implementation
6. **Testing** — unit and integration tests alongside source (`_test.go` and `_integration_test.go`)
7. **Mock Generation** — `//go:generate` comments; mock lives with whoever's tests consume it (see
   Mock Ownership above)
8. **Module Bootstrap** — dependencies wired in the module's own `module.go`
9. **Router Registration** — the module registers its own routes via `RegisterHttp`, invoked by
   `cmd/app`

## Beyond `internal/<module>/`

Not everything lives inside a single module. A few directories sit outside this
per-module structure — see their own `CLAUDE.md` for full rules:

- **`internal/platform/`** — shared I/O adapters (mailer, event bus, storage,
  HTTP error translation). Use when you need to talk to something *outside*
  the process and more than one module needs it. Never put business logic here.
- **`internal/shared/`** — shared domain-agnostic types with no I/O (errors,
  pagination, money). Use only for types every module would otherwise have to
  reinvent identically. When in doubt, duplicate in your module's `domain/`
  instead of reaching for `shared/` — bounded contexts should diverge by default.
- **`internal/bootstrap/`** — the composition root's module list
  (`Modules() []module.Module`, dependency-ordered). Imports every module, so it
  can live in neither `shared/` (must stay a leaf) nor `platform/` (no module
  imports). Consumed only by `cmd/app` (server) and `cmd/migrate` (migrations) so
  a new module is wired into both from one place. Wiring only — no logic.

**Rule of thumb:** if it's specific to your module's business rules, it stays in
`internal/<module>/`. If it's infrastructure with no business meaning, it's a
`internal/platform/` candidate. If it's a pure data shape with zero I/O and zero business
meaning, it's a `internal/shared/` candidate. If you're unsure which, default to keeping
it inside your own module — it's cheaper to extract later than to walk back
a bad shared dependency.
