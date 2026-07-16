# Scenario

**Feature**: `mark-done` forces log success with a required operator note (no SQL)

```
# operator marks a migration done without re-running SQL
seed failed row + matching fixture file
cli.Run(cfg, ["mark-done", id, "--note", note])
  -> logrepo.MarkDone -> status=success + note
  -> optional follow-up status shows skip
```

## Preconditions

- Positional migration_id must exist in log (seeded).
- `--note` required and non-empty after trim.
- Empty/missing note is usage exit 2 (no MarkDone call).
- DSN from Config for DB leaves.

## Steps

1. Happy leaf seeds failed row + fixture; runs mark-done; follow-up status.
2. empty-note leaf passes `--note ""` and expects exit 2.

## Context

- mark-done never Exec's migration SQL.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Happy leaf fills Config+MySQL itself; empty-note stays offline.
	req.CloseStdin = false
	return nil
}
```
