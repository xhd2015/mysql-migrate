## Expected

- `err == nil`, row found.
- `status == "success"`.
- `Note == req.Note` (non-empty).

## Errors

- None.

## Side Effects

- Row forced to success with operator note for audit.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("MarkDone: unexpected error: %v", err)
	}
	row := requireRow(t, resp)
	assertStatus(t, row, "success")
	if row.Note != req.Note {
		t.Fatalf("Note: got %q want %q", row.Note, req.Note)
	}
	assertNonEmpty(t, "Note", row.Note)
}
```
