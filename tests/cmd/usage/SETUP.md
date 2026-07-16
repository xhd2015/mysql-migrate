# Scenario

**Feature**: binary usage errors when DSN absent → nil cfg.DB (offline)

```
# DB subcommand without DSN (flag and env both absent) → usage exit 2
mysql-migrate --dir <tmp> status|apply
  -> no sql.Open; cfg.DB nil -> Error missing DB -> exit 2
```

## Preconditions

- Leaves under this branch expect exit **2** (usage).
- Migrate env stripped so ambient DSN cannot satisfy the binary.
- Offline — failure is config/usage, not DB connectivity.
- Split by subcommand: `apply-missing-dsn` vs `status-missing-dsn`.

## Steps

1. Leaf sets Args that omit DSN (and keep ClearMigrateEnv).
2. Exec binary; Assert exit 2 + missing/DB/dsn Error tokens.

## Context

- Sibling of `help/` (exit 0) and `status/` / `apply/` (happy DB path via Wrap).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Usage branch: force stripped env so ambient DSN cannot satisfy the binary.
	req.ClearMigrateEnv = true
	req.Env = nil
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
