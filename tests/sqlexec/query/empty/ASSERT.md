## Expected

- No harness error; `OpErr` nil.
- `QueryEmpty` is true.
- `QueryCount == 0`.

## Side Effects

- Temporary empty table created and dropped.

## Errors

- None expected (empty set is not an error).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.OpErr != nil {
		t.Fatalf("query_empty failed: %v", resp.OpErr)
	}
	if !resp.QueryEmpty {
		t.Fatal("QueryEmpty=false, want true (no Next)")
	}
	if resp.QueryCount != 0 {
		t.Fatalf("QueryCount=%d want 0", resp.QueryCount)
	}
}
```
