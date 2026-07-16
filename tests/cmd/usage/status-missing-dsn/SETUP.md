# Scenario

**Feature**: `status` without `--dsn` and without env DSN exits usage **2**

```
# dir provided so failure is missing DB (no edge open), not migrations dir
mysql-migrate --dir <tmp> status
  -> cfg.DB nil -> exit 2; Error mentions missing / DB / dsn
```

## Preconditions

- Args: `["--dir", <tempMigrationsDir>, "status"]` (no `--dsn`).
- `ClearMigrateEnv=true` so `MIGRATE_MYSQL_DSN` is not inherited.
- Temp dir exists (may be empty); failure must be missing DSN/DB, not "dir missing".
- Offline — no MySQL required.

## Steps

1. Create temp migrations dir.
2. Run status with only `--dir`.
3. Expect exit 2; combined output mentions missing/DB/dsn.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	req.MigrationsDir = dir
	req.ClearMigrateEnv = true
	req.Env = nil
	req.Args = []string{"--dir", dir, "status"}
	return nil
}
```
