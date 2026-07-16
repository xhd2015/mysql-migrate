## Expected

- No harness error; `OpErr` nil.
- `ScanValue` equals the seeded value (`77`).

## Side Effects

- Temporary table created and dropped.

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
		t.Fatalf("query_row_ok failed: %v", resp.OpErr)
	}
	want := int64(77)
	if len(req.SeedValues) > 0 {
		want = req.SeedValues[0]
	}
	if resp.ScanValue != want {
		t.Fatalf("ScanValue=%d want %d", resp.ScanValue, want)
	}
}
```
