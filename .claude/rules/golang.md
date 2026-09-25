---
paths:
  - "**/*.go"
---

# Go Code Rules

## Before Commit

Run `make unit-tests` and fix all issues.

## General Principles

1. **Follow existing patterns** - Check similar implementations in the target package before starting
2. **Use installed libraries** - Check `go.mod` before adding dependencies
3. **No sensitive data exposure** - Not even in logs
4. **Minimal comments** - Code should be self-documenting. Only add comments for non-obvious *why*. See [Comments](#comments) below for the full list of banned patterns.

## Comments

Write comments that are succinct and helpful. Focus on *why*, not *what*.

**Good comments:**
- Explain non-obvious business logic or constraints
- Document public API contracts (godoc style)
- Clarify complex algorithms or edge cases
- External constraints — why something must be done a specific way due to API/library quirks

**Exported identifier comments:**
- Do not add comments just because an identifier is exported if the only content is a restatement, e.g. `QuoteID is the quote ID`.
- Do not comment every struct field by default. Comment fields only when there is a non-obvious invariant, serialization contract, unit, lifecycle rule, nil/zero-value meaning, or external API requirement.
- Prefer one concise type-level comment over many field-level comments when the struct's purpose is clear from the field names.

**Bad comments:**
- Restate what the code does (`// increment counter` before `counter++`)
- Restate struct fields (`// QuoteID is the quote UUID`, `// CreatedAt is when it was created`)
- Explain constructor wiring that is obvious from the parameters and return value
- Explain the whole test scenario before a test whose name already says it
- Add multi-line "boundary tested" or QA-ticket blocks instead of using a precise test name
- Reference legacy code, previous implementations, or what was removed
- Include TODOs without ticket references if there is one
- Mention "for now", "temporary", or "will be replaced"
- Apologetic comments ("this is a hack", "sorry for this")

```go
// Good: explains why
// MarketplaceOrders skip payment validation because platforms handle payments externally.
if order.Channel == "Platform" {
    return nil
}

// Bad: references legacy
// Previously used the old EventService, now uses Publisher directly.
publisher.Publish(ctx, subject, payload)
```

### Comment Review Rule

**Do not flag comments for addition or removal** unless they:
- Contain secrets, credentials, or sensitive data (security issue)
- Reference removed/legacy code that no longer exists (misleading)
- Are factually incorrect and will confuse future readers
- Are structurally required by other rules
- Add broad PR-review narration or verbose test-story blocks that obscure the code path under review

The cost of flip-flopping on small comment preferences exceeds the benefit.
Accept existing comment style and move on unless the comment is noisy enough to
make the diff harder to review or maintain.

## Error Handling

### Domain Layer

Domain layer **MUST** return `errors.New` or `errors.Wrap` for all error cases, with a sentinel `*ErrorCode` defined via `errors.Define(...)`. Never use stdlib `errors`/`fmt.Errorf` in `domain/` — every returned error must be `*errors.Error` (this is what makes the handler assertion below safe).

```go
// Good - structured domain errors
var ErrInvalidCodeType = errors.Define("booking.invalid_code_type")

return errors.New(ErrInvalidCodeType, "invalid code type", errors.WithMetadata("codeType", t))

// Bad - raw errors
return fmt.Errorf("cart not found")  // Don't do this
return stderrors.New("invalid")      // Don't do this
```

Constructors: `internal/shared/errors/error.go`. Sentinel codes live per-module in `domain/` (e.g. `internal/identity/domain/user.go`).

### Application Layer

Application layer must wrap lower errors through `internal/shared/errors` too — `errors.Wrap(SentinelCode, err, "message")` — never pass a raw `error` up to the handler unwrapped:

```go
// Good
return errors.Wrap(domain.ErrBookingInternal, err, "we got a problem creating booking")

// Bad - raw error reaches the handler, breaks the errors.Is/type-assertion pattern below
return err
```

### API Handlers

Handlers are generated oapi-codegen strict-server methods (`func (s *Server) Op(ctx context.Context, request OpRequestObject) (OpResponseObject, error)`, see `internal/identity/api/http/server.go`) — no `http.ResponseWriter`/`response.*` helpers, no manual status codes. Check specific domain error sentinels with `errors.Is` and return the matching typed `OpNNNJSONResponse`, wrapping a small helper that builds the generated error body. The `err.(*errors.Error).Message()` call inside those helpers is safe *only* because Domain/Application layers guarantee `*errors.Error` — don't copy this pattern for errors that didn't come through the use case's `Handle`.

```go
func badRequest(err error) BadRequestJSONResponse {
    return BadRequestJSONResponse{Code: errors.BadRequestCode.String(), Message: err.(*errors.Error).Message()}
}

func unauthorized(err error) UnauthorizedJSONResponse {
    return UnauthorizedJSONResponse{Code: errors.Unauthorized.String(), Message: err.(*errors.Error).Message()}
}

func internalError(err error) InternalErrorJSONResponse {
    return InternalErrorJSONResponse{Code: errors.InternalCode.String(), Message: err.Error()}
}

func (s *Server) Login(ctx context.Context, request LoginRequestObject) (LoginResponseObject, error) {
    res, err := s.login.Handle(ctx, &query.LoginQuery{...})
    if err != nil {
        if errors.Is(err, domain.ErrAuthInternal) || errors.Is(err, domain.ErrAccountDeleted) {
            return Login401JSONResponse{unauthorized(err)}, nil
        }
        return Login500JSONResponse{internalError(err)}, nil
    }

    return Login200JSONResponse(res), nil
}
```

The error itself is never returned as the function's second value for a handled domain error — the sentinel maps to a typed response object instead, and `nil` is returned as the error. Only return a non-nil `error` for something oapi-codegen's own runtime should treat as unhandled.

#### Authenticated Principal Context

`middleware.Authenticate` validates a Bearer access token (signature, expiry, `type`, `iss`) and
puts a `Principal{Subject, Claims}` in context — it has no notion of a user. The module maps
`principal.Subject` (the `sub` claim) onto its own identifier. Rejected requests get a JSON 401
whose message carries only `crypto`'s sentinel errors, never the raw jwt library error.

```go
principal, ok := middleware.GetPrincipalFromContext(ctx)
if !ok {
    return OpNNN500JSONResponse{InternalErrorJSONResponse{Code: errors.InternalCode.String(), Message: "principal not found in context"}}, nil
}
```

### Error Codes Reference

Only four codes exist (`internal/shared/errors/code.go`). `Code.HttpStatus()` maps a code to its HTTP status.

| Code             | HTTP | Constant                | Use Case                |
| ---------------- | ---- | ------------------------ | ------------------------ |
| `INTERNAL_ERROR` | 500  | `errors.InternalCode`    | Server errors             |
| `BAD_REQUEST`    | 400  | `errors.BadRequestCode`  | Invalid input              |
| `NOT_FOUND`      | 404  | `errors.NotFoundCode`    | Resource doesn't exist    |
| `UNAUTHORIZED`   | 401  | `errors.Unauthorized`    | Auth required              |

### Response Construction (oapi-codegen strict server)

No response-writing helper package — the generated `OpNNNJSONResponse` types *are* the response contract, and their `VisitOpResponse(w)` method (generated) writes status + body. A handler's job is only to pick the right typed response, never to touch `http.ResponseWriter` directly:

```go
return Login200JSONResponse(res), nil        // 200, body = res
return Login400JSONResponse{badRequest(err)}, nil   // 400
return Login401JSONResponse{unauthorized(err)}, nil // 401
return Login500JSONResponse{internalError(err)}, nil // 500
```

Need to set something outside the JSON body (a cookie, a header)? Wrap the generated response type and override its `VisitOpResponse` method — don't reach for a response-writing helper.

There is no `DomainError` type, no `IsNotFoundError`/`IsValidationError`-style helpers, and no field-level validation error constructor. Check specific sentinel errors with `errors.Is(err, domain.ErrX)` against the `*ErrorCode` values each module defines via `errors.Define(...)`.

## Logging

Use `.claude/rules/logging.md` as the source of truth for all logging behavior.

## Review Guidelines

**Trust the compiler:** If code compiles successfully, do not flag:
- "Missing" methods or types (they may be from newer Go versions)
- Syntax that "doesn't exist" in the standard library
- API usage patterns - the compiler already validated these

**Go version:** This project uses Go 1.26+. Features like `sync.WaitGroup.Go()` are valid.

**Deliberate placeholders:** Mock/stub implementations (especially for payment, external APIs) are intentional until real implementations are needed. Do not flag these as issues unless they're in production code paths.
