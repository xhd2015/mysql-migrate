## Expected

- No harness error (failure is in Response.OpErr).
- `OpErr` is non-nil (invalid SQL rejected by driver/server).

## Side Effects

- None lasting (no successful DDL).

## Errors

- Expected: Exec error for invalid SQL.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.OpErr == nil {
		t.Fatal("expected non-nil OpErr for invalid SQL Exec")
	}
}
```
