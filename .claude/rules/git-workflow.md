# Git Workflow Rules

## Branch Naming

Format: `type/short-description`

| Type       | Use for                                   |
| ---------- | ----------------------------------------- |
| `feat`     | New user-facing features, API endpoints   |
| `fix`      | Bug fixes, regression patches             |
| `refactor` | Internal improvements, no behavior change |
| `docs`     | Documentation, README, guides             |
| `test`     | Test additions/improvements               |
| `chore`    | Build/tooling, dependencies               |

Examples: `feat/spot-api`, `fix/tax-entity-field`

## Commit Message Format

Format: `type(scope): description`

Examples:

- `feat(spot): implement spot validation`
- `fix(booking): avoid create booking twice`
- `refactor(order): extract validation logic`

**IMPORTANT: Never add Claude/AI co-author trailers** (e.g. `Co-Authored-By: Claude...`) to commit messages. Omit them entirely.

## Before Commit

**Never create a commit without passing tests.** If It fails, fix before committing.

```bash
make unit-tests
```

**ALWAYS review diff for:**
- API keys
- Credentials and tokens
- Passwords and private keys
- Database connection strings with passwords

**If secrets detected: STOP immediately.**

1. Stage the fix separately
2. Run `git reset` on the problematic commit
3. Remove secret from all commits if needed
4. Warn the team

## Secrets Management

### Local Development

| File | Purpose | Git Status |
|------|---------|------------|
| `.env` | Personal local overrides | Never commit (git-ignored) |
| `.env.example` | Template for setup | Committed |

### Production

Credentials are stored using Fly.io Secrets. Secrets are injected into the application as environment variables at runtime. The CI/CD pipeline does not access or retrieve application secrets.

## No Sensitive Data in Logs

**Never log:**
- Customer PII: email, phone, address, order history
- Payment information: card numbers (PCI-DSS), CVV, bank accounts
- Authentication tokens: JWT, session IDs, API keys
- Internal system credentials: database URLs with passwords, private keys

**Safe to log:**
- Card brand (Visa, Mastercard)
- Last 4 digits of card number
- Order IDs, store IDs (non-PII identifiers)
