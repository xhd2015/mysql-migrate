# Scenario

**Feature**: README states the mysql-migrate purpose (MySQL migrations)

```
# purpose section
README.md -> mentions "mysql-migrate" and migration / MySQL wording
```

## Preconditions

- README.md at module root (Classic RED if missing).
- Purpose prose may be free-form; required tokens are locked below.

## Steps

1. Set `req.Label` to `purpose`.
2. Require phrases: tool name, MySQL, and migrations (plural).

## Context

- Locks that the README is about this operator tool, not an empty placeholder.
- Does not require a fixed one-liner — only key tokens.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Label = "purpose"
	req.RequiredPhrases = []string{
		"mysql-migrate",
		"MySQL",
		"migration",
	}
	return nil
}
```
