# Scenario

**Feature**: EXACTLY-ONCE failed blocks; later pending deferred

```
# EO migration failed once — never auto re-apply; chain stops
files [EO-failed, later-pending]
logs [failed on EO]
  -> [blocked EO, deferred later], HasBlock=true
```

## Preconditions

- First file ExactlyOnce=true, log status `failed`.
- Second file non-EO, no log (would be apply if chain clear).

## Steps

1. Build EO file + later pending file.
2. Log only failed on EO.
3. Expect blocked then deferred; HasBlock=true.

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	const (
		idEO = "2026-07-16-01-[EXACTLY-ONCE]-drop-legacy"
		idB  = "2026-07-16-02-create-b"
	)
	req.Files = []inventory.MigrationFile{
		mf(idEO, true, "hash-eo"),
		mf(idB, false, "hash-b"),
	}
	req.Logs = []plan.LogRow{
		logRow(idEO, "failed", "hash-eo"),
	}
	return nil
}
```
