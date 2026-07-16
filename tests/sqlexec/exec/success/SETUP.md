# Scenario

**Feature**: Exec CREATE + INSERT reports RowsAffected >= 1

```
# happy Exec
DROP IF EXISTS -> CREATE TABLE -> INSERT one row
  -> Result.RowsAffected() >= 1
```

## Preconditions

- Isolated table name for this leaf.
- Op=`exec_ok`.

## Steps

1. Set Table via tableName("exok").
2. Set Op=`exec_ok`.
3. Assert ExecOK and RowsAffected >= 1.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "exec_ok"
	req.Table = tableName("exok")
	t.Logf("leaf exec/success: table=%s", req.Table)
	return nil
}
```
