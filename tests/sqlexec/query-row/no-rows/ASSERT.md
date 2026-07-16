## Expected

- No harness error.
- `OpErr` is non-nil.
- `OpErrIsNoRows` is true (`errors.Is(err, sql.ErrNoRows)`).

## Side Effects

- Temporary empty table created and dropped.

## Errors

- Expected: `sql.ErrNoRows` (or wrapped equivalent).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.OpErr == nil {
		t.Fatal("expected OpErr for QueryRow with no matching rows")
	}
	if !resp.OpErrIsNoRows {
		t.Fatalf("OpErrIsNoRows=false; OpErr=%v (want errors.Is sql.ErrNoRows)", resp.OpErr)
	}
}
```
