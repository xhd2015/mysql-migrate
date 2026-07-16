# Scenario

**Feature**: README lists every operator subcommand

```
# subcommands section
README.md -> status, plan, apply, mark-done, mark-failed, note, allow-retry
```

## Preconditions

- README.md at module root (Classic RED if missing).
- Subcommand set matches the sealed CLI surface (`tests/cli` / `tests/cmd`).

## Steps

1. Set `req.Label` to `subcommands`.
2. Require all seven subcommand names as substrings.

## Context

- Recovery commands (`mark-done`, `mark-failed`, `note`, `allow-retry`) must
  appear so operators discover the human gate, not only apply/status/plan.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Label = "subcommands"
	req.RequiredPhrases = []string{
		"status",
		"plan",
		"apply",
		"mark-done",
		"mark-failed",
		"note",
		"allow-retry",
	}
	return nil
}
```
