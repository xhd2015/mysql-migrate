# Scenario

**Feature**: binary status hand-off with `--dsn` + `--dir` (MySQL via Wrap)

```
# flags → edge open → sqlexec.Wrap → cfg.DB → cli.Run status
mysql-migrate --dsn $DSN --dir <fixture> status
  -> EnsureTable (+ optional ensured print) → plan table → exit 0
```

## Preconditions

- Leaves need reachable MySQL (harness DSN); **skip** when ping fails.
- Args must include global `--dsn` and `--dir` before `status`.
- MySQL leaves call `acquireMySQLExclusive` so ensure-created drop is safe.
- Unique fixture migration ids per session when files are written.

## Steps

1. Leaf ensures MySQL, acquires exclusive lock, writes fixtures if needed.
2. Sets Args with `--dsn`, `--dir`, `status` and AssertDSN for side effects.
3. Assert exit 0 and status table / ensure line as per leaf.

## Context

- Proves P2 edge path for status (not only apply). Full status matrix lives in `tests/cli/status/`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Status branch requires reachable MySQL; skip early if down.
	ensureMySQL(t, d)
	acquireMySQLExclusive(t)
	req.ClearMigrateEnv = true
	return nil
}
```
