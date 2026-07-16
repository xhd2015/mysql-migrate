# Scenario

**Feature**: success with ContentSHA256 mismatch → blocked + HashMismatch

```
# file bytes changed after a prior success — never silent re-apply same id
file hash=new, log success hash=old
  -> blocked, HashMismatch=true, reason hash_mismatch
  later pending -> deferred
```

## Preconditions

- First file success log with **different** ContentSHA256 than file.
- Second file pending.
- Locked: never re-apply same id on hash mismatch → **blocked**.

## Steps

1. File hash `hash-new`; log success with `hash-old`.
2. Later pending file.
3. Expect blocked + HashMismatch; later deferred.

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
		mf(idA, false, "hash-new"),
		mf(idB, false, "hash-b"),
	}
	req.Logs = []plan.LogRow{
		logRow(idA, "success", "hash-old"),
	}
	return nil
}
```
