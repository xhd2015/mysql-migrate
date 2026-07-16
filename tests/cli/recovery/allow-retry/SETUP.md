# Scenario

**Feature**: `allow-retry` clears EO failed for re-apply; rejects non-EO

```
cli.Run(cfg, ["allow-retry", id, "--note", note])
  -> EO failed: pending + note, exit 0
  -> non-EO: Error exit 1
```

## Preconditions

- DB leaf family.

## Steps

1. Leaves seed EO/non-EO failed and run allow-retry (optional follow-up apply).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fillConfigForDB(t, req)
	req.CloseStdin = false
	return nil
}
```
