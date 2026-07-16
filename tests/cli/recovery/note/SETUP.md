# Scenario

**Feature**: `note` updates operator note without changing status

```
cli.Run(cfg, ["note", id, "--note", note]) -> note updated; status unchanged
```

## Preconditions

- DB leaf family.

## Steps

1. Leaves seed + run note.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fillConfigForDB(t, req)
	req.CloseStdin = false
	return nil
}
```
