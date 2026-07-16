# Scenario

**Feature**: non-EXACTLY-ONCE failed still blocks the chain (pure plan)

```
# failed is blocked until human allow-retry; later pending deferred
files [failed-non-EO, later]
logs [failed]
  -> [blocked, deferred], HasBlock=true
```

## Preconditions

- First file ExactlyOnce=false, log status `failed`.
- Second file no log.

## Steps

1. Non-EO failed + later pending.
2. Expect blocked then deferred (not re-apply automatically in pure plan).

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	const (
		idA = "2026-07-16-01-create-a"
		idB = "2026-07-16-02-create-b"
	)
	req.Files = []inventory.MigrationFile{
		mf(idA, false, "hash-a"),
		mf(idB, false, "hash-b"),
	}
	req.Logs = []plan.LogRow{
		logRow(idA, "failed", "hash-a"),
	}
	return nil
}
```
