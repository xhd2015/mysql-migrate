# Scenario

**Feature**: binary apply hand-off with `--dsn` + `--dir` (MySQL)

```
# flags populate Config; cli.Run applies migrations
mysql-migrate --dsn $DSN --dir <fixture> apply
  -> plan clear → MarkRunning → Exec → MarkSuccess → exit 0
```

## Preconditions

- Leaves need reachable MySQL (harness DSN); **skip** when ping fails.
- Unique fixture migration ids / table names per session for isolation.
- Args must include global `--dsn` and `--dir` before `apply`.
- `multiStatements=true` on DSN for multi-statement files (single CREATE is fine).

## Steps

1. Leaf ensures MySQL, writes fixture SQL, cleans prior log/table.
2. Sets Args with `--dsn`, `--dir`, `apply` and AssertDSN for side effects.
3. Assert exit 0, progress, log success, table exists.

## Context

- Optional end-to-end wire check; full apply matrix lives in `tests/cli/apply/`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Apply branch requires reachable MySQL; skip early if down.
	// Leaves still write fixtures and set --dsn/--dir Args.
	ensureMySQL(t)
	req.ClearMigrateEnv = true
	return nil
}
```
