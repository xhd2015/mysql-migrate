# Scenario

**Feature**: `sqlexec` exports the locked interface types and `Wrap`

```
# package surface (offline — no MySQL)
import migrate/sqlexec
  -> type DB, Result, Rows, Row
  -> func Wrap(*sql.DB) DB
  -> method sets match locked API
```

## Preconditions

- Offline: no MySQL required.
- Package path: `github.com/xhd2015/mysql-migrate/migrate/sqlexec`.
- Classic RED: missing package → compile failure until implementer lands it.

## Steps

1. Set Op to `interface` for the method-set probe leaf.
2. Run reflection / compile-time probes in root `Run`.
3. Assert all types and required methods are present.

## Context

- Sibling of live MySQL branches and config branches.
- Pure API contract — no network, no side effects.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Interface branch: offline API surface only.
	req.Op = "interface"
	t.Log("interface branch: probe DB/Result/Rows/Row/Wrap (no MySQL)")
	return nil
}
```
