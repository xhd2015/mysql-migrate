## Expected

- `err == nil`.
- `resp.EnsureCallsOK` is true (both EnsureTable calls succeeded).
- `resp.TableExists` is true (`t_sql_migration_log` in current schema).

## Errors

- None on a healthy local MySQL.

## Side Effects

- May CREATE `t_sql_migration_log` if missing; second call must not error.
- Does not drop or truncate existing data.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("EnsureTable twice: unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if !resp.EnsureCallsOK {
		t.Fatal("EnsureCallsOK: got false want true")
	}
	if !resp.TableExists {
		t.Fatal("TableExists: t_sql_migration_log missing after EnsureTable")
	}
}
```
