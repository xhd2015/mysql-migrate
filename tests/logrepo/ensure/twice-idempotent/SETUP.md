# Scenario

**Feature**: EnsureTable twice succeeds and table is present

```
# two EnsureTable calls on local DB
EnsureTable -> EnsureTable -> information_schema has t_sql_migration_log
```

## Preconditions

- Local MySQL reachable.
- No migration_id required.

## Steps

1. Keep Op=ensure_twice from parent.
2. Expect both EnsureTable calls nil and TableExists=true.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Explicit leaf contract: two EnsureTable calls, no row key needed.
	req.Op = "ensure_twice"
	req.MigrationID = "" // schema-only scenario
	t.Logf("ensure twice-idempotent: Op=%s (MySQL ensured by root Setup)", req.Op)
	return nil
}
```
