# Scenario

**Feature**: method sets for DB, Result, Rows, Row match the locked API

```
# compile + reflect method names
sqlexec.DB     -> Exec, Query, QueryRow, Close
sqlexec.Result -> LastInsertId, RowsAffected
sqlexec.Rows   -> Next, Scan, Close, Err
sqlexec.Row    -> Scan
sqlexec.Wrap   -> present as func
```

## Preconditions

- Op remains `interface` from parent (leaf may reaffirm).
- No table or DSN needed.

## Steps

1. Reaffirm Op=`interface`.
2. Run probes; expect all Has* and *MethodsOK flags true.

## Context

- First leaf implementers fix when creating the package skeleton.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "interface"
	t.Log("leaf methods-present: full interface method set")
	return nil
}
```
