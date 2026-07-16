## Expected

- No harness error; non-nil response.
- `WrapNonNil` is true.
- `ExecOK` is true (SELECT 1 succeeded) and `ScanValue == 1`.
- `OpErr` is nil.

## Side Effects

- Read-only `SELECT 1` on the harness database.

## Errors

- None expected.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.WrapNonNil {
		t.Fatal("Wrap returned nil DB")
	}
	if resp.OpErr != nil {
		t.Fatalf("SELECT 1 via Wrap failed: %v", resp.OpErr)
	}
	if !resp.ExecOK || resp.ScanValue != 1 {
		t.Fatalf("SELECT 1: ExecOK=%v ScanValue=%d want 1", resp.ExecOK, resp.ScanValue)
	}
}
```
