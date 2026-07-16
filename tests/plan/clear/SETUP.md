# Scenario

**Feature**: clear plan chain — only skip/apply, never blocked or deferred

```
# classifications are success→skip or pending→apply only
plan.Build -> HasBlock=false
  items Action ∈ {skip, apply}
```

## Preconditions

- Every file classifies as **skip** or **apply** (no failed/unknown/running/mismatch).
- Stop-chain never engages: no deferred items.

## Steps

1. Leaf provides files + matching success or absent/pending logs.
2. Assert all actions are skip/apply and `HasBlock=false`.

## Context

- Sibling of `blocked/` where stop-chain and human intervention apply.

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	// Clear-chain leaves must only produce skip/apply (HasBlock=false).
	if req.Files == nil {
		req.Files = []inventory.MigrationFile{}
	}
	if req.Logs == nil {
		req.Logs = []plan.LogRow{}
	}
	return nil
}
```
