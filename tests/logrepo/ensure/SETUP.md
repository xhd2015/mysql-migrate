# Scenario

**Feature**: EnsureTable creates `t_sql_migration_log` idempotently

```
# CREATE TABLE IF NOT EXISTS — safe to call repeatedly
logrepo.EnsureTable(db) -> ok
logrepo.EnsureTable(db) -> ok again (no error, no drop)
```

## Preconditions

- `req.Op` is `ensure_twice` for all descendants.
- No prior row assertions; only schema existence matters.

## Steps

1. Set `req.Op = "ensure_twice"`.
2. Run calls EnsureTable twice and checks information_schema.
3. Assert table exists and both calls succeeded.

## Context

- First call may create the table on a fresh DB; second must still succeed.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "ensure_twice"
	return nil
}
```
