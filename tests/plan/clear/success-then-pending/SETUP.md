# Scenario

**Feature**: success then pending → skip then apply

```
# first applied successfully; second never run
files [A,B]
logs [success A matching hash]
  -> [skip A, apply B], HasBlock=false
```

## Preconditions

- File A has success log with matching hash.
- File B has no log row.

## Steps

1. Two files in order A then B.
2. Only A in logs as success.
3. Expect skip, apply.

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
	)
	req.Files = []inventory.MigrationFile{
		mf(idA, false, hA),
		mf(idB, false, "hash-b"),
	}
	req.Logs = []plan.LogRow{
		logRow(idA, "success", hA),
	}
	return nil
}
```
