# Platform Folder Structure

Adapters to the outside world (I/O, third-party services) shared across modules.
If it makes a network call, touches disk, or talks to infra outside our process — it lives here (`internal/platform`).

## Admission Test
A package belongs in `internal/platform/` only if it answers YES to:
1. Does it perform real I/O (network, disk, external API)?
2. Is it free of business/domain logic — no knowledge of Order, Payment, User, etc.?
3. Would this same code make sense in a completely unrelated project?

If any answer is NO, it does not belong here.

## Forbidden
- No imports from any `internal/<module>/domain` or `internal/<module>/app`.
- No business rules, no module-specific types (e.g. `platform/mailer` must not
  have an `OrderConfirmationEmail` type — that template/type lives in the
  `order` module; mailer only exposes `Send(to, subject, body string) error`).
- No direct handling of a specific module's error kinds beyond the shared AppError contract.

## Consumption Rule
Each `internal/platform/<pkg>` package owns both its interface and its implementation
together (e.g. `eventbus.EventPublisher` interface + `eventbus.Publisher` struct
live in the same package).

- Only `internal/<module>/app` may import a `internal/platform/<pkg>`
  **interface type** directly (e.g. `eventbus.EventPublisher`, `mailer.EmailSender`).
- `internal/<module>/domain` must NEVER import anything from `internal/platform`.
  Domain stays free of any technical/infrastructure dependency, including interfaces.
- `internal/<module>/app` must NEVER construct or reference a platform
  **concrete type** (e.g. `&eventbus.Publisher{}`). Concrete types are only
  constructed in the module's own `module.go` (its `Init()`, per
  [folder-structure.md](folder-structure.md)) and injected as the interface via
  the use case's constructor.
- `internal/<module>/module.go` and `internal/<module>/api` may import 
  and construct `internal/platform` concrete types freely 
  (this is expected — e.g. wiring, health checks, admin tooling).

## Examples packages
- `internal/platform/mailer` — SMTP/SendGrid client
- `internal/platform/eventbus` — pub/sub for domain events
- `internal/platform/httpkit` — AppError → HTTP response translation, middleware
- `internal/platform/storage` — blob/S3 client

## Adding a new package
Requires: (1) it passes the admission test above, (2) it's referenced by 2+ modules.
A single-consumer package stays constructed inside that module's own `module.go`
(alongside its other platform concretes) unless the team explicitly agrees to promote it
to `internal/platform`.
