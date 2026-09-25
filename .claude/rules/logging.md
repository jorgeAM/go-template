---
paths:
  - "**/*.go"
---

# Logging Rules

## Source of Truth

This file is the canonical logging rule for the repository.

Use structured JSON logging everywhere.

## Approved APIs

### App Code

Use the `internal/platform/log` package (`github.com/jorgeAM/tuky-api/internal/platform/log`), backed by `ZapLogger`.

All app-code log calls take `context.Context` first, then a stable message, then fields built from `log.With*` option constructors:

- `log.Debug(ctx, "message", log.WithString("key", value))`
- `log.Info(ctx, "message", log.WithString("key", value))`
- `log.Warn(ctx, "message", log.WithString("key", value))`
- `log.Error(ctx, "message", log.WithError(err))`

Available field constructors: `WithString`, `WithInt`, `WithInt32`, `WithInt64`, `WithFloat32`, `WithFloat64`, `WithBool`, `WithDuration`, `WithTime`, `WithBinary`, `WithError`, `WithStack`, `WithObject` (see Hot-Path Hardening for `WithObject` limits).

If no `context.Context` is available at the call site, obtain or thread one (e.g. `context.Background()` at bootstrap) rather than introducing a parallel no-context API — `internal/platform/log` has one call shape per level.

## Required Format

Every log entry must follow this shape:

- message: stable event name or outcome
- fields: flat key-value pairs for dynamic values
- output: JSON

```go
log.Info(ctx, "booking created",
	log.WithString("booking_id", booking.ID),
)
```

```go
log.Error(ctx, "spot cannot be validated",
	log.WithString("spot_id", spot.ID),
)
```

## Required Conventions

- Use `snake_case` field names
- Use flat fields; avoid nested maps unless the payload is intentionally structured
- Log errors as `"error", err`
- Keep the message short and stable
- Put dynamic values in fields, not in the message

## Hot-Path Hardening

In high-frequency paths such as middleware, webhook handlers, retries or  polling loops:

- Prefer scalar fields such as IDs, counts, enums, booleans, durations, and status values
- Prefer compact summaries over full objects
- Log only the fields needed for debugging, operations, or auditability
- Treat `zap.Any(...)` as a fallback for small or rare values, not as the default shape for hot-path logs

Avoid logging full domain structs or payload-like objects in hot paths, including:

- booking
- spot
- request objects
- response objects
- webhook payloads
- event bus events (`events.Event`) — its `Payload` field is an untyped `interface{}` that can carry PII or full domain structs

If a payload must be logged, log a compact subset instead of the whole object.

```go
// Good
log.Info(ctx, "stripe webhook received",
	log.WithString("event_type", event.Type),
	log.WithString("order_id", event.OrderID),
)

// Bad
log.Info(ctx, "stripe webhook received",
	log.WithObject("payload", event),
)
```

```go
// Good — event listener/dispatch loop
log.Warn(ctx, "event handler not found",
	log.WithString("topic", event.Topic.String()),
	log.WithString("event_id", event.ID),
)

// Bad — dumps event.Payload (interface{}) via zap.Any
log.Warn(ctx, "event handler not found",
	log.WithString("topic", event.Topic.String()),
	log.WithObject("event", event),
)
```

## Banned Patterns

- `fmt.Sprintf(...)` inside log calls
- Placeholder formatting in log messages such as `%s`, `%v`, `%d`
- Embedding dynamic `key=value` data inside the message when it should be a field
- Logging secrets, tokens, JWTs, auth headers, raw credentials, or full request/response bodies
- Logging whole hot-path objects when a compact scalar subset would do

```go
// Bad
log.Info(ctx, "order %s created for store %s", order.ID, order.StoreNumber)

// Bad
log.Info(ctx, fmt.Sprintf("order %s created", order.ID))

// Bad
log.Info(ctx, "order created order_id=123 store_number=10145")
```
### Raw Zap

Do not use raw `*zap.Logger` in app code, domain, or application layers. Always go through `internal/platform/log`.

The only permitted raw-zap usage is inside `internal/platform/log` itself (`zap.go`, `default.go`) to build and configure the underlying `ZapLogger`. That code is the logging package's own implementation, not a general-purpose escape hatch.

If a third-party API genuinely requires a `*zap.Logger` value, add an accessor/adapter to `internal/platform/log` (e.g. exposing the underlying logger) rather than instantiating zap directly at the call site.


## Data Preservation

When refactoring logging:

- do not remove log events without explicit justification
- preserve level and lifecycle placement
- preserve all previously logged values as structured fields
- normalize message wording only if the event meaning stays intact

## Sensitive Data

Never log:

- passwords
- tokens
- JWTs
- session secrets
- card data
- raw credential material
- full request or response bodies unless explicitly approved

Be especially careful with email, phone, address, and other PII. Log only when required for operations or auditability.
