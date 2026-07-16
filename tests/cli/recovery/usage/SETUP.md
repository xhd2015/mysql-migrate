# Scenario

**Feature**: recovery flag-shape usage errors (missing note / id)

```
# offline parse validation — no --local/--remote
mark-done without --note or without migration_id -> exit 2
```

## Preconditions

- Offline leaves; no MySQL required.
- Representative command: `mark-done`.

## Steps

1. Leaves set incomplete Args.
2. Assert exit 2 and message tokens.

## Context

- Replaces lifelog missing-target with note/id only (Config carries DSN).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.CloseStdin = false
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
