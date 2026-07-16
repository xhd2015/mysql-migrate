# Scenario

**Feature**: mark-done with empty `--note` is a usage error (exit 2)

```
# --note is required and must be non-empty after trim
cli.Run(cfg, ["mark-done", "any-id", "--note", ""])
  -> exit 2; Error about --note / missing note
```

## Preconditions

- Args include `--note` with empty string value.
- Offline — no DB seed required (fails at parse before open).
- DSN may be empty; parse should fail on note first or still exit 2.

## Steps

1. Set mark-done with empty note.
2. Expect exit 2; combined output mentions note.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Synthetic id; parse fails before DB. No fillConfigForDB.
	req.Args = []string{"mark-done", "p5-empty-note-id", "--note", ""}
	return nil
}
```
