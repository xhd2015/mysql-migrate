# Scenario

**Feature**: recovery without migration_id is usage exit 2

```
cli.Run(cfg, ["mark-done", "--note", "x"])  # no positional id
  -> exit 2; message about missing migration_id
```

## Preconditions

- --note present; no positional migration_id.
- Offline.

## Steps

1. Set Args without id.
2. Expect exit 2; output mentions migration_id or missing.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"mark-done", "--note", "p5 missing id note"}
	return nil
}
```
