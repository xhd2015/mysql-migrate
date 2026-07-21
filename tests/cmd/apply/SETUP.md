# Scenario

**Feature**: binary apply hand-off with `--dsn` + `--dir` (MySQL via Wrap)

```
# flags → edge open → sqlexec.Wrap → cfg.DB → cli.Run apply
mysql-migrate --dsn $DSN --dir <fixture> apply
  -> plan clear → MarkRunning → Exec → MarkSuccess → exit 0
```

## Preconditions

- Leaves need reachable MySQL (harness DSN); **skip** when ping fails.
- Unique fixture migration ids / table names per session for isolation.
- Args must include global `--dsn` and `--dir` before `apply`.
- `multiStatements=true` on DSN for multi-statement files (single CREATE is fine).
- MySQL leaves call `acquireMySQLExclusive` so ensure-created drop is safe.

## Steps

1. Leaf ensures MySQL, acquires exclusive lock, writes fixture SQL, cleans prior log/table.
2. Sets Args with `--dsn`, `--dir`, `apply` and AssertDSN for side effects.
3. Assert exit 0, progress, log success, table exists.

## Context

- Optional end-to-end wire check; full apply matrix lives in `tests/cli/apply/`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Apply branch requires reachable MySQL; skip early if down.
	// Exclusive lock serializes with status/ensure-created DROP.
	ensureMySQL(t, d)
	acquireMySQLExclusive(t)
	req.ClearMigrateEnv = true
	return nil
}
```
