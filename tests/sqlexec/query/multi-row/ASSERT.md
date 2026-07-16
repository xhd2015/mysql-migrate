## Expected

- No harness error; `OpErr` nil.
- `QueryCount == 2`.
- `Scanned` equals `req.SeedValues` (`[10, 20]`).

## Side Effects

- Temporary table created and dropped via test cleanup.

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
		t.Fatalf("query_multi failed: %v", resp.OpErr)
	}
	if resp.QueryCount != len(req.SeedValues) {
		t.Fatalf("QueryCount=%d want %d; scanned=%v", resp.QueryCount, len(req.SeedValues), resp.Scanned)
	}
	if len(resp.Scanned) != len(req.SeedValues) {
		t.Fatalf("Scanned=%v want %v", resp.Scanned, req.SeedValues)
	}
	for i := range req.SeedValues {
		if resp.Scanned[i] != req.SeedValues[i] {
			t.Fatalf("Scanned[%d]=%d want %d (full=%v)", i, resp.Scanned[i], req.SeedValues[i], resp.Scanned)
		}
	}
}
```
