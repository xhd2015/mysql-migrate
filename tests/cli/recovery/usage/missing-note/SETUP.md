# Scenario

**Feature**: recovery without `--note` flag is usage exit 2

```
cli.Run(cfg, ["mark-done", "some-id"])  # no --note
  -> exit 2; message about missing --note
```

## Preconditions

- migration_id present; `--note` absent entirely.
- Offline.

## Steps

1. Set Args without --note.
2. Expect exit 2; output mentions note.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"mark-done", "p5-missing-note-id"}
	return nil
}
```
