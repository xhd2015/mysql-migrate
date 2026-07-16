# Scenario

**Feature**: after a blocked item, later success remains skip (not deferred)

```
# stop-chain only rewrites apply → deferred; skip is unchanged
files [failed, success-matching, pending]
  -> [blocked, skip, deferred]
```

## Preconditions

- First file failed (blocked).
- Second file success with matching hash (would skip even without block).
- Third file no log (would apply → deferred after block).

## Steps

1. Three files: failed, success, pending.
2. Assert middle stays **skip**; last **deferred**.

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
		idC = "2026-07-16-03-create-c"
		hB  = "hash-b"
	)
	req.Files = []inventory.MigrationFile{
		mf(idA, false, "hash-a"),
		mf(idB, false, hB),
		mf(idC, false, "hash-c"),
	}
	req.Logs = []plan.LogRow{
		logRow(idA, "failed", "hash-a"),
		logRow(idB, "success", hB),
	}
	return nil
}
```
