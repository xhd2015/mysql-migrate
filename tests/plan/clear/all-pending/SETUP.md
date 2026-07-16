# Scenario

**Feature**: all files pending (no log rows) → all apply

```
# two files, no logs
files [A, B], logs []
  -> [apply A, apply B], HasBlock=false
```

## Preconditions

- Two non-EXACTLY-ONCE files with distinct IDs, sorted by FileName.
- No log rows.

## Steps

1. Construct two `MigrationFile` entries in filename order.
2. Leave `Logs` empty.
3. Expect both actions **apply**, `HasBlock=false`.

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	req.Files = []inventory.MigrationFile{
		mf("2026-07-16-01-create-a", false, "aaa111"),
		mf("2026-07-16-02-create-b", false, "bbb222"),
	}
	req.Logs = []plan.LogRow{}
	return nil
}
```
