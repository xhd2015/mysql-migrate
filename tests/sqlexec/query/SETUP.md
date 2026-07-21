# Scenario

**Feature**: `DB.Query` returns a Rows cursor for multi-row reads

```
# multi-row path
DB.Query(ctx, SELECT ...) -> Rows
  -> Next/Scan until done; Close; Err
```

## Preconditions

- Live MySQL required.
- Leaves seed isolated tables.

## Steps

1. ensureMySQL.
2. Leaf chooses multi-row vs empty outcome.

## Context

- Outcome siblings: multi-row vs empty result set.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureMySQL(t, d)
	t.Log("query branch: live MySQL for Query multi/empty")
	return nil
}
```
