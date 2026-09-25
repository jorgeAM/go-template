---
paths:
  - "internal/*/domain/*.go"
---

# Domain Layer Rules

Rules for entities/aggregates inside `internal/<module>/domain/`. Complements
[folder-structure.md](folder-structure.md) (where files go) — this file governs
how the domain code inside those files must behave.

## Scope

Applies to every module's `internal/<module>/domain/`. Every aggregate in the
codebase already complies (private fields, getters, invariant-enforcing
mutation methods) — there is no standing grandfathered exception. New
aggregates must comply from the start. Constructor input structs
(`NewXInput`) are not the aggregate and are expected to have public fields —
don't confuse the two when checking compliance. If a genuinely non-compliant
aggregate ever surfaces, fix it opportunistically when touched for another
reason rather than carving out a permanent exception for it.

## Rules

1. **Private fields, accessed only through getters.** No exported fields on
   an entity/aggregate struct. Go's unexported-field rule already blocks
   `application`/`infrastructure` from reaching in directly — this rule is
   about not exporting fields on the struct in the first place.

2. **Mutation happens through methods — avoid anemic models.** No setter that
   just assigns a field. A mutation method enforces whatever invariant makes
   that mutation valid, and is named for the domain action (`AttachAvatar`,
   `UpdateLanguage`), not `SetAvatar`/`SetLanguage`.

3. **Mutation methods that can violate an invariant must return `error`, not
   silently no-op.** If a mutation can't apply, the caller must be told —
   don't swallow the failure.

4. **Repository interface lives in its own file, next to the aggregate it
   serves** — `domain/[name]_repository.go`, not inlined into the entity
   file. This is folder-structure.md's existing convention
   (`domain/[name]_repository.go`); this rule doesn't restate it, just
   flags that the entity file and the repository file are siblings, not one
   file wearing two hats.

5. **One repository per aggregate, not one per table.** The repository
   interface exposes aggregate-root operations only (`Save`, `FindByID`,
   `FindByEmail`, ...). Child entities inside the aggregate (e.g.
   `BookingItem` inside `Booking`) are persisted through the aggregate
   root's repository — they never get their own repository.

6. **Collections stay encapsulated.** Don't expose an aggregate's internal
   slice/map as a public field or a getter that returns the live reference.
   Return a copy, or expose read access through a method; mutate only
   through aggregate methods (`AddItem`, not `booking.Items = append(...)`).

7. **The persistence-reconstruction constructor lives in its own
   `domain/unmarshall.go`, not in the entity file.** Private fields mean the
   aggregate's own `New*` constructor re-runs creation invariants — wrong for
   rebuilding from a DB row, which must skip them. Give that path a separate
   `Unmarshall<Entity>(...)` function (see `domain.UnmarshallUser` in
   `internal/identity/domain/unmarshall.go`) that only the module's own `adapters/db` package
   calls. One `unmarshall.go` per module, one function per aggregate inside
   it — not a new file per aggregate.

8. **If an entity/aggregate turns out to have no behavior (anemic model),
   move it down to `app/models/`** as a plain struct, rather than forcing
   domain ceremony (private fields, getters, an interface-only repository)
   onto something that's genuinely just a data holder. Promote it back to
   `domain/` once it starts acquiring real behavior/invariants — don't
   design for that ahead of time.

   **The deciding test is invariant enforcement, not mutation methods or
   persistence.** An aggregate with no post-construction mutation methods is
   not automatically anemic — if its constructor enforces real invariants
   (validation, required relationships) and it carries its own identity,
   that is domain behavior, just front-loaded into construction instead of
   spread across setters. This is normal for create-once, never-edited
   aggregates. Likewise, having a repository does not by itself rescue an
   entity from being anemic — persistence is an infrastructure necessity
   every aggregate root needs, not evidence of business logic. Ask: does the
   constructor enforce anything a raw struct literal couldn't? If no, and
   there are no mutation methods either, it's anemic — move it.

   **When it moves, its repository interface moves with it**, out of
   `domain/` and into `app/models/<name>_repository.go` — colocated with the
   entity struct it serves, same reasoning as rule 4 above (repository lives
   next to what it repositories), just relocated to `models/` since that's
   where the entity itself now lives. The postgres implementation in
   `adapters/db/` is unaffected beyond its import path. See
   [folder-structure.md](folder-structure.md) for the full layout.

## Value Objects

Every module may define its own value objects inside its `domain/` folder
(e.g. `order.Amount`, `order.Quantity`, `booking.Schedule`) — these are
in scope for this rule and already follow rules 1-2 by convention: private
fields, validated in a `New*` constructor, immutable — mutation replaces the
whole value, there is no partial in-place mutation. This rule extends the
same shape to entities/aggregates.

Shared value objects in `internal/shared/valueobject/` (`valueobject.Email`,
`valueobject.UUID`, ...) are out of this rule's path scope — they're governed by
[shared-folder-structure.md](shared-folder-structure.md) instead. They
happen to follow the same private-field/`New*`/immutable shape, but that's
convention, not this rule reaching into `shared/`.
