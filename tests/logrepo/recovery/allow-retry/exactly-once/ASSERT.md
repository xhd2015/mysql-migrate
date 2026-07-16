## Expected

- `err == nil`, row found.
- `status == "pending"` (clear for plan re-apply).
- `ExactlyOnce == true` (flag preserved).
- `Note == req.Note`.

## Errors

- None.

## Side Effects

- EO failed row becomes pending so plan.Build can apply again; audit note kept.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("AllowRetry EO: unexpected error: %v", err)
	}
	row := requireRow(t, resp)
	assertStatus(t, row, "pending")
	if !row.ExactlyOnce {
		t.Fatal("ExactlyOnce: got false want true")
	}
	if row.Note != req.Note {
		t.Fatalf("Note: got %q want %q", row.Note, req.Note)
	}
}
```
