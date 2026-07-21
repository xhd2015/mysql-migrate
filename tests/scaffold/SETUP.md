# Scenario

**Feature**: empty repo becomes a buildable module with public migrate.Config (DB-only)

```
# caller imports migrate package and constructs Config
caller -> import github.com/xhd2015/mysql-migrate/migrate
caller -> migrate.Config{DB, MigrationsDir, ProgramName, AppliedBy}

# module builds with optional stub packages
module root (go.mod) -> go build ./... -> exit 0
```

## Preconditions

- Working directory for the product is the mysql-migrate repo root
  (`d.DOCTEST_ROOT/../..` from this tree).
- Module path in `go.mod` is `github.com/xhd2015/mysql-migrate` (already present).
- Go toolchain (`go`) is on PATH.
- Package `github.com/xhd2015/mysql-migrate/migrate` with type `Config` is the
  production surface under test (Classic RED until Config.DB lands).
- Config fields: `DB` (`sqlexec.DB`), `MigrationsDir`, `ProgramName`,
  `AppliedBy`. **No `DSN` field**.
- Out of scope: inventory/plan/logrepo/cli behavior, MySQL, SQL files.
  Live `sqlexec.Wrap` coverage lives under `tests/sqlexec/`.

## Steps

1. Root setup only validates the request pointer; leaves set `req.Mode` and
   any sample Config field values.
2. Root `Run` switches on `req.Mode` to construct `Config`, reflect the type,
   or run `go build ./...` at the module root.

## Context

- Config package: prefer `github.com/xhd2015/mysql-migrate/migrate`
  (`migrate/config.go` or `migrate/migrate.go`).
- Locked fields: `DB` (`sqlexec.DB`), `MigrationsDir`, `ProgramName`,
  `AppliedBy`.
- Stub packages may be empty `doc.go` / `package` files so `go build ./...` works.

```go
import (
	"fmt"
	"os/exec"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	if req == nil {
		return fmt.Errorf("nil request")
	}
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go not found in PATH: %w", err)
	}
	t.Log("scaffold root: Config.DB-only contract (no DSN field)")
	return nil
}
```
