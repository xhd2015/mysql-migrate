# Scenario

**Feature**: empty files and empty logs → empty plan, HasBlock=false

```
# both slices empty
plan.Build([], []) -> {Items: [], HasBlock: false}
```

## Preconditions

- `Files` is empty.
- `Logs` is empty.

## Steps

1. Set both input slices to empty.
2. Run `plan.Build`.

```go
import (
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

func Setup(t *testing.T, req *Request) error {
	req.Files = []inventory.MigrationFile{}
	req.Logs = []plan.LogRow{}
	return nil
}
```
