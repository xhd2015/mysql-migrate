# Scenario

**Feature**: blocked plan chain — HasBlock=true and stop-chain defers later apply

```
# first blocked item stops the apply chain
classify → blocked sets seenBlock
later would-be apply → deferred
later skip stays skip
HasBlock = true whenever any item is blocked
```

## Preconditions

- At least one file classifies as **blocked** (failed, unknown, stale running, hash mismatch).
- Later pending files become **deferred**, not apply.

## Steps

1. Leaf constructs files+logs that trigger a specific block reason.
2. Assert blocked item, deferred successors, and `HasBlock=true`.

## Context

- Human intervention required before apply can proceed past the block.

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	// Blocked-chain leaves must produce at least one blocked item (HasBlock=true).
	if req.Files == nil {
		req.Files = []inventory.MigrationFile{}
	}
	if req.Logs == nil {
		req.Logs = []plan.LogRow{}
	}
	return nil
}
```
