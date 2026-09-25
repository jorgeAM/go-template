---
name: create-pr
description: Create a pull request using the PR template. Use when ready to open a PR for your changes.
---

# Create Pull Request

Automates PR creation and proper template filling.

## Arguments

| Argument | Description |
|----------|-------------|
| `--draft` | Create as draft PR (default) |
| `--ready` | Create as ready for review |

## Steps

### 1. Pre-Flight Checks

Verify environment before proceeding:

```bash
# Check gh CLI is available
which gh || echo "MISSING"

# Check current branch (must not be main/master/dev directly)
BRANCH=$(git branch --show-current)
echo "Branch: $BRANCH"

# Determine base branch:
# - if current branch is "dev" -> base is "main" (release PR)
# - else if "dev" exists on remote -> base is "dev" (normal feature PR)
# - else fall back to repo's configured default branch
if [ "$BRANCH" = "dev" ]; then
  BASE="main"
elif git show-ref --verify --quiet refs/remotes/origin/dev; then
  BASE="dev"
else
  BASE=$(gh repo view --json defaultBranchRef -q .defaultBranchRef.name)
fi
echo "Base: $BASE"

# Check for commits ahead of base
git log origin/$BASE..HEAD --oneline

# Check for existing PR (scoped to this head branch, not just current default)
gh pr list --head "$BRANCH" --state open --json number,url,baseRefName 2>&1 || echo "NO_EXISTING_PR"

# Check for uncommitted changes
git status --porcelain
```

**Handle failures:**

| Check | Failure Response |
|-------|-----------------|
| `gh` missing | Show manual instructions: "Install with `brew install gh` and run `gh auth login`" |
| On main/master | STOP: "Create a feature branch first: `git checkout -b type/ISSUE-ID/description`" |
| No commits | STOP: "No commits to create PR from. Commit your changes first." |
| Existing PR | If list is non-empty, show URL(s) + base branch(es) and ask: "PR already exists. Open it?" (if multiple bases, ask which one) |
| Uncommitted changes | WARN: "You have uncommitted changes. Continue anyway?" |

**Note:** `dev` is a normal working branch here (not a protected/release-only branch) — the "On main/master" guard above does not apply to `dev`. Feature branches are cut from and PR'd into `dev`; only a PR *from* `dev` itself targets `main` (release).

### 2. Extract PR Information

Gather information from git:

```bash
# Get branch name for PR title construction
BRANCH=$(git branch --show-current)
echo "Branch: $BRANCH"

# Get commit messages for summary (uses $BASE from step 1)
git log origin/$BASE..HEAD --pretty=format:"%s" --reverse

# Get first commit for PR title (usually describes the feature)
git log origin/$BASE..HEAD --pretty=format:"%s" --reverse | head -1

# Get changed files for context
git diff origin/$BASE --name-only

# Get diff stats
git diff origin/$BASE --stat

# Get per-file churn for reviewer notes
git diff origin/$BASE --numstat
```

**Extract from branch name:**
- Pattern: `type/slug` (type values: see Branch Naming in `.claude/rules/git-workflow.md`)
- Type: first segment, e.g. `feat/spot-api` → `feat`

**Build PR title:**
- Format: same as commit message format, see `.claude/rules/git-workflow.md`
- Example: `feat(cart): add checkout validation`
- Single commit following the format: reuse it as-is
- Multiple commits: synthesize a title summarizing the overall change across all commits (type/scope from branch + changed files, description covers the whole diff, not just the first commit)

### 3. Fill PR Template

Read and fill `.github/PULL_REQUEST_TEMPLATE.md` verbatim. **CRITICAL: Preserve ALL template sections, HTML comments, and `<details>` blocks exactly.**

**Template sections to fill:**

| Section | How to fill |
|---------|-------------|
| **Summary** | 2-4 bullets summarizing changes from commit messages |
| **Reviewer Notes** | Help reviewers focus by grouping relevant files into business logic vs config/generated/low-risk. Include changed-file count context, low-risk-heavy areas, and highest-signal files to review. |
| **Test Plan** | Ask user OR derive from changes (e.g., "Unit tests added for X"); always include `make test` as verification step |
| **Checklist** | Mark `[x]` for: self-reviewed (always), tests (if tests changed), `make test` passes (if run), no secrets committed (always) |
| **Schema / interface changes (`<details>`)** | Evaluate each of the 3 items independently — do not check/uncheck as a block: (1) check migration item only if `database/migration/` changed; (2) check mocks item only if any `port/`/`domain/service/` interface changed AND `mocks/` was regenerated in the diff; (3) check breaking-changes item only if breaking changes are explicitly written out in Context below — if there are none, leave it unchecked (not N/A-checked) |
| **Context** | Leave HTML comment placeholder unless user provides context (related issues, incidents, design notes) |

**Reviewer Notes requirements:**

- Start with a short changed-file count/churn note, for example: `This PR changes 8 files, most churn in the transaction module and its tests.`
- Include one bullet for **Business logic / highest-signal review**. Name concrete files/directories owning behavior, e.g. `internal/booking/application/...`, `internal/transaction/domain/...`, handler files.
- Include one bullet for **Config / generated / low-risk**. Name concrete paths such as `database/migration/...`, `mocks/...`, `config/dependencies.go` wiring.
- If tests are a large part of the diff, call them out separately rather than letting their volume obscure the business-logic review targets.
- Do not use vague labels like "misc files" — every category must list representative paths.

**Example filled template:**

```markdown
## Summary

- Add cart validation before checkout
- Return validation errors with field-level details
- Add unit tests for validation logic

## Reviewer Notes

- This PR changes 9 files, most churn in tests and mock updates.
- Business logic / highest-signal review: `internal/cart/domain/validation.go`, `internal/cart/handler.go`.
- Config / generated / low-risk: `internal/cart/mocks/...`, `database/migration/000015_add_cart_constraints.up.sql`.

## Test Plan

- Unit tests cover validation rules
- Run `make test` to verify

## Checklist

- [x] Self-reviewed the diff
- [x] Tests cover new/changed behavior
- [x] `make test` passes locally
- [x] No secrets or credentials committed

<details>
<summary>Schema / interface changes? Expand for additional items</summary>

- [x] Migration added via `make new_migration` (up + down)
- [ ] Mocks regenerated via `make generate` if interfaces changed
- [ ] Breaking changes called out explicitly in Context below (or N/A)

</details>

## Context

<!-- Why is this change needed? Link related issues, incidents, design notes. -->
```

### 4. Create PR

Push branch and create PR:

```bash
# Ensure branch is pushed with tracking
git push -u origin $(git branch --show-current)

# Create PR against the base determined in step 1 (draft by default)
gh pr create \
  --base "$BASE" \
  --title "feat(cart): add checkout validation" \
  --body "$(cat <<'EOF'
<filled template content>
EOF
)" \
  --draft
```

**Arguments:**
- Default: `--draft`
- If `--ready` passed: omit `--draft`

### 5. Success Output

After successful PR creation, display:

```
PR created successfully!

  Title: feat(cart): add checkout validation
  URL: https://github.com/jorgeAM/tuky-api/pull/153
  Status: Draft

Next steps:
  - Wait for CI checks to pass
  - Mark ready when done: gh pr ready
  - Address review feedback
```

## Error Handling

| Error | Response |
|-------|----------|
| `gh pr create` fails | Show error, suggest: "Check `gh auth status`" |
| Tests fail | Show output, ask to fix and retry |
| Push fails | Show error, suggest: "Check remote access with `git remote -v`" |

## Usage

Invoke with:
```
/create-pr          # Create draft PR (default)
/create-pr --ready  # Create ready-for-review PR
```
Or ask: "create a PR", "open PR for my changes", "submit this for review"
