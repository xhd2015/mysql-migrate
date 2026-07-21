# Scenario

**Feature**: first DB cmd after missing log table prints ensure-created via Wrap path

```
# drop t_sql_migration_log under exclusive lock; status recreates via cli EnsureTable
DROP TABLE t_sql_migration_log
mysql-migrate --dsn <harness> --dir <empty> status
  -> stdout: ensured: t_sql_migration_log (created)
  -> exit 0
```

## Preconditions

- MySQL reachable; **skip** when down.
- `acquireMySQLExclusive` held (parent) before DROP so parallel cmd leaves wait.
- Empty migrations dir (status still runs EnsureTable then empty plan).
- Do **not** call harness EnsureTable after DROP (binary must observe created=true).

## Steps

1. Parent: ensureMySQL + exclusive lock.
2. DROP `t_sql_migration_log` if present.
3. Args: `--dsn` + `--dir` empty tmp + `status`.
4. Expect exit 0 and exact ensure-created line on stdout.

## Context

- Proves edge Wrap path still feeds `logrepo.EnsureTable` → cli print when created.
- Table is recreated by the binary; later leaves see it again after unlock.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Parent already ensureMySQL + acquireMySQLExclusive.
	dropMigrationLogTable(t, d)

	dir := t.TempDir()
	dsn := harnessDSN()
	req.MigrationsDir = dir
	req.AssertDSN = dsn
	req.ClearMigrateEnv = true
	req.Args = []string{"--dsn", dsn, "--dir", dir, "status"}
	return nil
}
```
