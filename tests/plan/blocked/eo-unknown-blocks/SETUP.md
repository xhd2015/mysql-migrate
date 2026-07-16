# Scenario

**Feature**: EXACTLY-ONCE with status unknown is blocked

```
# EO + unknown requires human; later pending deferred
files [EO, later]
logs [unknown on EO]
  -> [blocked, deferred], HasBlock=true
```

## Preconditions

- First file ExactlyOnce=true, log status `unknown`.
- Second file pending (no log).

## Steps

1. EO file + later file.
2. Unknown log on EO only.
3. Expect blocked, deferred, HasBlock.

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
		logRow(idEO, "unknown", "hash-eo"),
	}
	return nil
}
```
