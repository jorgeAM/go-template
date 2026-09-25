## Summary

<!-- 2-4 bullets: what changed, at a high level -->

-
-

## Reviewer Notes

<!-- Help the reviewer prioritize. Which files carry the actual logic vs. boilerplate? -->

- This PR changes N files, most churn in `...`.
- Business logic / highest-signal review: `path/to/file.go`, `path/to/other.go`.
- Config / generated / low-risk: `database/migration/...`, `mocks/...`.

## Test Plan

<!-- How did you verify this? -->

- Run `make unit-tests` to verify

## Checklist

- [ ] Self-reviewed the diff
- [ ] Tests cover new/changed behavior
- [ ] `make unit-tests` passes locally
- [ ] No secrets or credentials committed

<details>
<summary>Schema / interface changes? Expand for additional items</summary>

- [ ] Migration added via `make new_migration` (up + down)
- [ ] Mocks regenerated via `make generate` if interfaces changed
- [ ] Breaking changes called out explicitly in Context below (or N/A)

</details>

## Context

<!-- Why is this change needed? Link related issues, incidents, design notes. -->
