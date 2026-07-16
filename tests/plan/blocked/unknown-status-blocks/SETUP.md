# Scenario

**Feature**: status unknown (non-EO) is blocked

```
# unknown is always a human gate
file non-EO, log unknown
  -> blocked, HasBlock=true
```

## Preconditions

- One non-EXACTLY-ONCE file with log status `unknown`.
- Optional second pending to show deferral.

## Steps

1. Unknown on first; second pending.
2. Expect blocked, deferred.

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
		logRow(idA, "unknown", "hash-a"),
	}
	return nil
}
```
