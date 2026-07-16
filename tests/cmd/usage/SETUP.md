# Scenario

**Feature**: binary usage errors for incomplete global config (offline)

```
# apply without DSN (flag and env both absent) → usage exit 2
mysql-migrate --dir <tmp> apply
  -> Error missing DSN -> exit 2 (no MySQL open)
```

## Preconditions

- Leaves under this branch expect exit **2** (usage).
- Migrate env stripped so ambient DSN cannot satisfy the binary.
- Offline — failure is config/usage, not DB connectivity.

## Steps

1. Leaf sets Args that omit DSN (and keep ClearMigrateEnv).
2. Exec binary; Assert non-zero usage exit + dsn Error tokens.

## Context

- Sibling of `help/` (exit 0) and `apply/` (happy DB path).

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
