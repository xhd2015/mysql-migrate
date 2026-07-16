## Expected

- No harness error.
- `ExecOK` is true.
- `RowsAffected >= 1` for the INSERT.
- `OpErr` is nil.

## Side Effects

- Temporary table created and dropped (cleanup in Run).

## Errors

- None expected.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.OpErr != nil {
		t.Fatalf("exec_ok failed: %v", resp.OpErr)
	}
	if !resp.ExecOK {
		t.Fatal("ExecOK=false, want true")
	}
	if resp.RowsAffected < 1 {
		t.Fatalf("RowsAffected=%d, want >= 1", resp.RowsAffected)
	}
}
```
