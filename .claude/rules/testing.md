---
paths:
  - "**/*_test.go"
  - "**/*_integration_test.go"
  - "internal/**/mocks/*"
---

# Testing Rules

## Mandatory Test Coverage

**These are requirements, not suggestions:**

| Code Changed | Required Test Type | Run Command |
|--------------|-------------------|-------------|
| Repository implementations | Integration test | `make integration-tests` |
| Business logic, calculations | Unit test | `make unit-tests` |

**Before completing any task:**
1. Write the required tests
2. Run and verify they pass
3. Do not mark work complete until tests pass

## Unit vs Integration Tests

**Dependencies determine test type, not the name:**

| Test Type | Uses | File Naming | Run Command |
|-----------|------|-------------|-------------|
| Unit | Mocks/stubs | `*_test.go` | `make unit-tests` |
| Integration | Testcontainers/real DB | `*_integration_test.go` | `make integration-tests` |

**Never name a mock-based test as `*_integration_test.go`.**

| Scenario | Test Type |
|----------|-----------|
| Business logic, handlers, calculations | Unit |
| Database queries, cross-service workflows | Integration |

## Required Patterns

### Time and Date Stability

Tests must not depend on fixed calendar dates that will age into a different
business meaning, such as a hard-coded "future" fulfillment date that eventually
becomes today or the past.

Use one of these instead:

- A frozen clock or injected `now` when the code supports it
- Relative dates derived from test-local `now`
- A named fixture date only when it is tied to seeded data and the test asserts that exact fixture behavior

Avoid mixing real wall-clock validation with hard-coded ISO timestamps unless
the timestamp is deliberately part of the fixture contract.

### Parallel Test Execution

Use `t.Parallel()` for unit tests that don't share state:

| Safe to parallelize | NOT safe to parallelize |
|---------------------|-------------------------|
| Fresh mocks per test | Shared package-level variables |
| Isolated structs | Global state modification |
| Table-driven tests with `t.Parallel()` in subtests | Tests that depend on execution order |

**Table-driven with parallel subtests:**

```go
func TestValidation(t *testing.T) {
    t.Parallel()
    tests := []struct{ name, input string }{...}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // Each subtest runs concurrently
            // ...
        })
    }
}
```

### Integration Tests

Integration tests **MUST** include `//go:build integration` and skip under `testing.Short()`:

```go
//go:build integration

package mypackage

func TestMyPackage_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    // use real DB (testcontainers), no mocks
}
```

Container setup pattern (pooling, shared lifecycle across tests) will be defined once the
first integration test is written — do not invent helper names ahead of that.

## Test Organization

**Unit tests:**
- Table-driven tests covering happy path, error cases, and edge conditions together
- Individual tests only when a case needs setup too divergent to tabulate
- Mocks generated via `gomock` (`go.uber.org/mock`), run `make generate`, living wherever the
  tests that consume them live (`internal/<module>/mocks/` for that module's own tests). A
  module's `client/` interface is mocked by each *consumer* module in the consumer's own
  `mocks/`, never by the producing module — see Mock Ownership in
  [folder-structure.md](folder-structure.md)

**Integration tests:**
- Real DB via `testcontainers`
- Organize by operation (`t.Run("Create")`, `t.Run("Get")`, `t.Run("Update")`, `t.Run("Delete")`, etc.)

