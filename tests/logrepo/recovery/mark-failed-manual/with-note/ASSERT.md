## Expected

- `err == nil`, row found.
- `status == "failed"`.
- `Note == req.Note`.

## Errors

- None.

## Side Effects

- Row forced to failed with operator note.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("MarkFailedManual: unexpected error: %v", err)
	}
	row := requireRow(t, resp)
	assertStatus(t, row, "failed")
	if row.Note != req.Note {
		t.Fatalf("Note: got %q want %q", row.Note, req.Note)
	}
}
```
