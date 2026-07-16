# Scenario

**Feature**: empty inputs produce an empty plan

```
# no migration files
plan.Build([], logs?) -> Plan{Items: [], HasBlock: false}
```

## Preconditions

- Zero migration files in the request.
- Logs may be empty or irrelevant (no files to classify).

## Steps

1. Leave `Files` empty (or explicitly empty slice).
2. Build plan and expect no items and no block.

## Context

- Edge of the status machine: nothing to apply or block.

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	// Empty branch starts with zero files; leaves may keep or replace.
	req.Files = []inventory.MigrationFile{}
	req.Logs = []plan.LogRow{}
	return nil
}
```
