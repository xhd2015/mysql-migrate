# Scenario

**Feature**: `DB.QueryRow` returns a single-row accessor

```
# single-row path
DB.QueryRow(ctx, SELECT ...).Scan(&dest) -> value | sql.ErrNoRows
```

## Preconditions

- Live MySQL required.

## Steps

1. ensureMySQL.
2. Leaf chooses one-row vs no-rows.

## Context

- Outcome siblings: one-row success vs no-rows error.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureMySQL(t, d)
	t.Log("query-row branch: live MySQL for QueryRow")
	return nil
}
```
