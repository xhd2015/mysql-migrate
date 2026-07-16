# Scenario

**Feature**: log status empty string or "pending" means apply (allow-retry cleared)

```
# after human allow-retry, status cleared to pending/empty
files [A status="", B status="pending"]
  -> [apply A, apply B], HasBlock=false
```

## Preconditions

- Two files each with a log row present (so not "missing"), but status is
  `""` or `"pending"` (post allow-retry).
- Prefer rule: empty string or `"pending"` → **apply**.

## Steps

1. File A log Status `""`; file B log Status `"pending"`.
2. Expect both **apply**.

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
		mf(idA, false, "h-a"),
		mf(idB, false, "h-b"),
	}
	req.Logs = []plan.LogRow{
		logRow(idA, "", "old-a"),
		logRow(idB, "pending", "old-b"),
	}
	return nil
}
```
