# Scenario

**Feature**: stale status running treated as unknown → blocked (EO and non-EO)

```
# crashed mid-apply leaves running forever → treat as unknown
files [non-EO running, EO running, later pending]
  -> both running items blocked with LogStatus=unknown
  -> later deferred
```

## Preconditions

- Two files with log status `running`: one ExactlyOnce=false, one ExactlyOnce=true.
- Third file pending (no log).
- Locked rule: any `running` → effective LogStatus **unknown** → **blocked**.

## Steps

1. Three files in FileName order: non-EO, EO, later.
2. Logs: running on first two.
3. Expect both blocked with LogStatus `unknown`; third deferred; HasBlock.

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	const (
		idA  = "2026-07-16-01-create-a"
		idEO = "2026-07-16-02-[EXACTLY-ONCE]-risky"
		idC  = "2026-07-16-03-create-c"
	)
	req.Files = []inventory.MigrationFile{
		mf(idA, false, "hash-a"),
		mf(idEO, true, "hash-eo"),
		mf(idC, false, "hash-c"),
	}
	req.Logs = []plan.LogRow{
		logRow(idA, "running", "hash-a"),
		logRow(idEO, "running", "hash-eo"),
	}
	return nil
}
```
