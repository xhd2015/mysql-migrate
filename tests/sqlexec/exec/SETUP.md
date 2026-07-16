# Scenario

**Feature**: `DB.Exec` runs statements with context and returns Result or error

```
# mutating path
DB.Exec(ctx, CREATE|INSERT|DROP, args...) -> Result{RowsAffected} | error
```

## Preconditions

- Live MySQL required.
- Leaves assign unique `Table` names via `tableName`.

## Steps

1. ensureMySQL.
2. Leaf chooses success vs bad-sql Op and table/query.

## Context

- Outcome siblings: success vs bad-sql.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureMySQL(t)
	t.Log("exec branch: live MySQL for Exec success/error")
	return nil
}
```
