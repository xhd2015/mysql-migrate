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
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureMySQL(t)
	t.Log("query-row branch: live MySQL for QueryRow")
	return nil
}
```
