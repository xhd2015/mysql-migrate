# Scenario

**Feature**: `sqlexec.Wrap` adapts `*sql.DB` into a usable `DB`

```
# harness opens DSN; Wrap produces facade
sql.Open(harnessDSN) -> sqlexec.Wrap(sqlDB) -> non-nil DB
DB.QueryRow(ctx, "SELECT 1").Scan(&n) -> n==1
```

## Preconditions

- Live MySQL required; skip when unreachable.
- Harness owns `sql.Open`; production code under test only uses `Wrap`.

## Steps

1. ensureMySQL.
2. Leaf sets Op=`wrap`.
3. Run wraps and probes with `SELECT 1`.

## Context

- Sibling of interface (offline) and exec/query (mutating) branches.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureMySQL(t)
	req.Op = "wrap"
	t.Log("wrap branch: require live MySQL for SELECT 1 via Wrap")
	return nil
}
```
