# Scenario

**Feature**: success with empty log ContentSHA256 still skips (no mismatch)

```
# success log without recorded hash is treated as compatible
file hash=abc, log status=success ContentSHA256=""
  -> skip (not hash_mismatch)
```

## Preconditions

- One file with non-empty ContentSHA256.
- Log status `success` and **empty** ContentSHA256.
- Locked rule: empty log hash does **not** count as mismatch → **skip**.

## Steps

1. Single file with hash `abc123`.
2. Success log with empty hash field.
3. Expect **skip**, HashMismatch=false.

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	const id = "2026-07-16-01-create-a"
	req.Files = []inventory.MigrationFile{
		mf(id, false, "abc123"),
	}
	req.Logs = []plan.LogRow{
		logRow(id, "success", ""), // empty log hash
	}
	return nil
}
```
