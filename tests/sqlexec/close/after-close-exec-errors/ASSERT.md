## Expected

- No harness error.
- `CloseErr` is nil (first Close succeeds).
- `PostCloseExecErr` is true (Exec after Close fails).

## Side Effects

- Underlying `*sql.DB` closed (harness must not reuse it).

## Errors

- Expected: non-nil error from post-close Exec.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.CloseErr != nil {
		t.Fatalf("Close() error: %v", resp.CloseErr)
	}
	if !resp.PostCloseExecErr {
		t.Fatal("expected Exec after Close to return an error")
	}
}
```
