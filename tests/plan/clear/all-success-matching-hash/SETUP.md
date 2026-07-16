# Scenario

**Feature**: all success with matching content hashes → all skip

```
# every file has success log with same ContentSHA256 as file
files [A,B] logs success(A,hashA) success(B,hashB)
  -> [skip A, skip B], HasBlock=false
```

## Preconditions

- Two files with known hashes.
- Log rows status=`success` and ContentSHA256 equals each file hash.

## Steps

1. Build two files with hashes `hash-a` / `hash-b`.
2. Add matching success log rows.
3. Expect both **skip**.

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
		hA  = "hash-a"
		hB  = "hash-b"
	)
	req.Files = []inventory.MigrationFile{
		mf(idA, false, hA),
		mf(idB, false, hB),
	}
	req.Logs = []plan.LogRow{
		logRow(idA, "success", hA),
		logRow(idB, "success", hB),
	}
	return nil
}
```
