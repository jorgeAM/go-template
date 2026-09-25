---
name: code-auditor
description: Run multi-agent code review locally before creating a PR. Use when you want to review changes before pushing.
---

# Local Code Review

Run a multi-agent code review locally against your branch diff. Catches bugs, rule violations, and missing test coverage before opening a PR.

## Review Agents

The following agents run in parallel to analyze your changes:

| Agent | Focus |
|-------|-------|
| `rule-keeper` | Audit changes against CLAUDE.md and `.claude/rules/*.md` for violations |
| `bug-hunter` | Find bugs, both obvious in the diff and ones requiring broader codebase context |
| `test-guardian` | Check if new/changed code has appropriate test coverage |
| `gopher` | go-template Go conventions — domain boundaries, error handling, testing patterns (only if `.go` files changed) |

## Steps

### 1. Run preflight checks

Run tests first to filter out issues tooling already catches (no linter configured yet — see `.claude/rules/golang.md` for `golangci-lint`/`gofmt` mentions, which agents should still defer to once one exists):

```bash
make unit-tests
```

If preflight fails, fix those issues first before proceeding with code review.

### 2. Get the diff to review

Determine the base branch the same way `create-pr` does (see `.claude/skills/create-pr/SKILL.md`): feature branches diff against `dev`; `dev` itself diffs against `main` (release).

```bash
BRANCH=$(git branch --show-current)
if [ "$BRANCH" = "dev" ]; then
  BASE="main"
elif git show-ref --verify --quiet refs/remotes/origin/dev; then
  BASE="dev"
else
  BASE=$(gh repo view --json defaultBranchRef -q .defaultBranchRef.name)
fi
git diff $BASE...HEAD
```

### 3. Run review agents in parallel

Use the Agent tool to launch all 4 agents (`rule-keeper`, `bug-hunter`, `test-guardian`, `gopher`) in a single message with multiple tool calls, run in the foreground since their output feeds validation next. Skip `gopher` if no `.go` files changed. Each agent should receive:
- The diff output
- Instructions to return findings as a JSON array

Each finding must include:
- `description`: What the issue is
- `file`: File path
- `line`: Line number (if applicable)
- `confidence`: 0-100 how confident the agent is
- `severity`: "critical", "important", or "minor"
- `reason`: Why this is an issue

Example Agent prompt for each agent:
```
Review this diff for [agent focus area]. Return findings as a JSON array.

Diff:
[paste diff here]

Return format:
[{"description": "...", "file": "...", "line": 42, "confidence": 90, "severity": "important", "reason": "..."}]
```

### 4. Validate issues (reduce false positives)

After all agents complete, validate each issue with confidence >= 80:

For each candidate issue, use a haiku subagent to verify:
```
Verify this issue is real by checking the code context.

Issue: [description]
File: [file]
Line: [line]

Read the file and determine:
1. Is this actually a problem? (not a false positive)
2. Is this HIGH SIGNAL? (see criteria below)

Return JSON: {"valid": true/false, "reason": "why valid or invalid"}
```

**HIGH SIGNAL criteria** (from Anthropic's code review plugin):
- Code that fails to compile or has unresolved references
- Clear logic errors that will produce wrong results
- Unambiguous CLAUDE.md/rule violations with exact quote

**NOT high signal** (exclude these):
- Style concerns or subjective improvements
- Potential issues dependent on specific inputs
- Issues linters will catch (golangci-lint, gofmt)
- "Might be a problem" concerns

### 5. Aggregate validated results

After validation:
1. Collect validated issues only
2. **Track raw counts per agent** (total found before validation)
3. **Track validated counts per agent** (kept vs rejected)
4. Group kept issues by severity

### 6. Output format

Display results in terminal:

**If no issues found:**
```markdown
## Code Auditor Results

✅ **No issues found.** Checked for bugs, logic errors, and CLAUDE.md compliance.

**Branch:** `feat/checkout-validation`
**Files changed:** 5

### What Went Well
[2-4 positive observations about the code]

### Recommendation

✅ **APPROVED** - Ready for PR creation.

---
Agents: rule-keeper, bug-hunter, test-guardian, gopher
```

**If issues found:**
```markdown
## Code Auditor Results

**Branch:** `feat/checkout-validation`
**Files changed:** 5

### Review Stats
| Agent | Found | Validated | Rejected |
|-------|-------|-----------|----------|
| rule-keeper | 4 | 2 | 2 |
| bug-hunter | 2 | 1 | 1 |
| test-guardian | 0 | 0 | 0 |
| gopher | 6 | 0 | 6 |
| **Total** | **12** | **3** | **9** |

### Critical Issues (X)

**path/to/file.go:42** - Nil pointer dereference
Agent: bug-hunter | Confidence: 95%

> user.Name accessed without nil check after GetUser() which can return nil

```suggestion
if user == nil {
    return errors.New("user not found")
}
```

### Important Issues (Y)
[issues with optional suggestions, or "None found"]

### Minor Issues (Z)
[issues or "None found"]

### What Went Well
Based on the diff, highlight 2-4 positive observations about the code:
- Good practices followed (error handling, naming, structure)
- Clean patterns used (separation of concerns, appropriate abstractions)
- Test coverage quality

### Recommendation

⚠️ **CHANGES REQUESTED** - Please address critical/important issues before PR.

or

✅ **APPROVED WITH SUGGESTIONS** - Minor issues noted but not blocking.

---
Agents: rule-keeper, bug-hunter, test-guardian, gopher
```

**Recommendation logic:**
- Any **critical** issues → CHANGES REQUESTED
- Only **important** issues → CHANGES REQUESTED
- Only **minor** issues → APPROVED WITH SUGGESTIONS
- No issues → APPROVED

**Committable suggestions:** For small, obvious fixes (1-5 lines), include a `suggestion` code block that shows the exact fix. This allows quick resolution without back-and-forth.

## When to Use

- Before creating a PR
- After completing a feature to catch issues early

## Review Philosophy

Based on [Anthropic's code review plugin](https://github.com/anthropics/claude-code/blob/main/plugins/code-review/commands/code-review.md):

**HIGH SIGNAL only.** Report only issues that are clearly problems:
- Code that won't compile or has unresolved references
- Logic errors that will produce wrong results
- Unambiguous rule violations with exact quotes

**Validation eliminates false positives.** Each candidate issue is verified by a subagent before reporting. This prevents the frustrating cycle of fixing issues only to have new ones appear on iteration 2-6.

**Comments are not flagged.** Per `.claude/rules/golang.md`, agents do not suggest adding or removing comments unless they contain security issues, reference removed/legacy code, or are factually incorrect.
