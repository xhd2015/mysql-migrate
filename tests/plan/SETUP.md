# Scenario

**Feature**: pure migration plan status machine — Build files+logs → actions

```
# plan is pure: files + log rows -> ordered actions (no I/O)
plan.Build(files, logs) -> Plan{Items, HasBlock}
  classify each file: skip | apply | blocked
  stop-chain: after first blocked, later apply -> deferred
```

## Preconditions

- Module: `github.com/xhd2015/mysql-migrate`.
- Target package (to implement):
  `migrate/plan`
  import `github.com/xhd2015/mysql-migrate/migrate/plan`
- Inventory types already exist:
  `github.com/xhd2015/mysql-migrate/migrate/inventory.MigrationFile`
- Rules locked for classification (per file, then stop-chain):
  - Source of truth for ExactlyOnce: **file.ExactlyOnce**.
  - No log / status `""` / `"pending"` → **apply**.
  - `success` + hash match **or** log hash empty → **skip**.
  - `success` + non-empty log hash ≠ file hash → **blocked**, HashMismatch=true, reason contains `hash_mismatch`.
  - `failed` / `unknown` → **blocked** (EO and non-EO).
  - `running` → effective LogStatus **`unknown`** → **blocked** (stale).
  - Stop-chain: after first **blocked**, later **apply** → **deferred**; **skip** stays **skip**.
  - `HasBlock` iff any item Action is **blocked**.
- Items order: **FileName ascending** (Build sorts if needed).
- No MySQL, no filesystem I/O in Build.

## Steps

1. Leaf Setup fills `req.Files` and `req.Logs` in memory.
2. `Run` calls `plan.Build` and maps to `Response.Items` / `HasBlock`.
3. Leaf Assert checks actions, LogStatus, HashMismatch, HasBlock.

## Context

- Classic RED: `plan.LogRow` / `plan.Build` missing until implementer ports API.
- Standalone root under `tests/plan/` — does not inherit inventory Request/Run.
- Pure in-memory: construct `MigrationFile` + `LogRow` in Setup (no FS).

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	if req.Files == nil {
		req.Files = []inventory.MigrationFile{}
	}
	if req.Logs == nil {
		req.Logs = []plan.LogRow{}
	}
	return nil
}

// mf builds a minimal MigrationFile for plan tests (ID + FileName + flags + hash).
func mf(id string, exactlyOnce bool, contentSHA string) inventory.MigrationFile {
	return inventory.MigrationFile{
		ID:            id,
		FileName:      id + ".sql",
		ExactlyOnce:   exactlyOnce,
		ContentSHA256: contentSHA,
	}
}

// logRow builds a plan.LogRow for the given migration id and status.
func logRow(id, status, contentSHA string) plan.LogRow {
	return plan.LogRow{
		MigrationID:   id,
		Status:        status,
		ContentSHA256: contentSHA,
	}
}
```
