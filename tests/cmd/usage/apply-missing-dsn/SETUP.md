# Scenario

**Feature**: `apply` without `--dsn` and without env DSN exits usage **2**

```
# dir provided so failure is DSN, not migrations dir
mysql-migrate --dir <tmp> apply
  -> exit 2; Error mentions dsn / missing
```

## Preconditions

- Args: `["--dir", <tempMigrationsDir>, "apply"]` (no `--dsn`).
- `ClearMigrateEnv=true` so `MIGRATE_MYSQL_DSN` is not inherited.
- Temp dir exists (may be empty); failure must be DSN, not "dir missing".
- Offline — no MySQL required.

## Steps

1. Create temp migrations dir.
2. Run apply with only `--dir`.
3. Expect exit 2; combined output mentions dsn/missing.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	req.MigrationsDir = dir
	req.ClearMigrateEnv = true
	req.Env = nil
	req.Args = []string{"--dir", dir, "apply"}
	return nil
}
```
