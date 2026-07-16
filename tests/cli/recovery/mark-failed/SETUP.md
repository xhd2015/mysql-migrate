# Scenario

**Feature**: `mark-failed` forces failed + operator note

```
cli.Run(cfg, ["mark-failed", id, "--note", note]) -> log failed + note
```

## Preconditions

- DB leaf family; DSN required.

## Steps

1. Leaves seed + run mark-failed.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fillConfigForDB(t, req)
	req.CloseStdin = false
	return nil
}
```
