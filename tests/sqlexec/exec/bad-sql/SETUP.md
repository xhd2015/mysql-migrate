# Scenario

**Feature**: Exec of invalid SQL returns a non-nil error

```
# error path
DB.Exec(ctx, "THIS IS NOT VALID SQL !!!") -> error != nil
```

## Preconditions

- Live MySQL (server must parse/reject SQL).
- Op=`exec_err`.
- Query is a clearly invalid statement.

## Steps

1. Set Op=`exec_err` and a fixed bad Query string.
2. Assert OpErr is non-nil.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "exec_err"
	req.Query = "THIS IS NOT VALID SQL !!!"
	t.Log("leaf exec/bad-sql: expect driver/server error")
	return nil
}
```
