# Shared Folder Structure

## Purpose
Domain-agnostic types with zero I/O and zero business rules, reused across modules
to avoid duplication or illegal cross-module domain imports.

## Admission Test
A type belongs in `internal/shared/` only if it answers YES to:
1. Is it free of any single module's business rules (no Order, Payment, User semantics)?
2. Would every module's domain layer independently need to invent this if it didn't exist?
3. Is it stable — unlikely to change because one module's requirements change?

## Forbidden
- No type named after, or shaped around, a specific bounded context
  (no `PaymentAmount` here — `Money` is fine, `PaymentAmount` is not).
- No repository interfaces, no use cases — those are module-owned.
- No dependencies on any `internal/<module>/*` package (this must stay a leaf
  dependency; if it needs to import a module, it's not shared, it's misplaced).
- No dependencies on `internal/platform/*` either. If a type needs I/O, it's not shared.
- Standard library only, plus this third-party allowlist:
  - `github.com/go-chi/chi/v5` — the router in the `module.Module` contract
  - `github.com/google/uuid` — `valueobject.UUID` (v7 generation, parsing)
  - `github.com/golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt` — `crypto` primitives

  Adding a library to this list requires explicit team agreement, same as a new package.

## Consumption Rule
`internal/shared/*` may be imported from any layer of any module —
`internal/<module>/domain`, `internal/<module>/app`, `internal/<module>/api`, or
`internal/<module>/adapters` — since it imports nothing internal beyond
itself and carries no business or I/O concerns. This is the one exception to
"domain imports nothing technical": shared types aren't technical, they're
just common vocabulary every module speaks.

## Example package
- `internal/shared/errors/` — `Error`, `ErrorCode`, `Define`, `Wrap`, `WithMetadata`
- `internal/shared/criteria/` — filters, ordering and pagination
- `internal/shared/valueobject/email.go` — Email value object
- `internal/shared/valueobject/uuid.go` — UUID value object

One package per concept (`errors`, `criteria`, `valueobject`), not one
package per type — a new value object is a new file inside `valueobject/`,
not a new sub-package, unless it grows large enough to need its own
internal helpers.

## When Not to Share
If two modules need "the same" type but their business meaning differs even
slightly (e.g. Order's "Money" has currency rounding rules Payment doesn't),
DO NOT force it into shared — duplicate it in each module's own `domain/`
instead. Bounded contexts diverging is a feature, not a bug.

## Adding a new package
Requires: (1) it passes the admission test above, (2) at least 2 modules
independently need the exact same shape today (not "might need it later" —
speculative sharing is how this folder rots into a dumping ground).
