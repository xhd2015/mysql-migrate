## Expected

- `err == nil`, row found.
- `status == "success"` (unchanged from seed).
- `Note == req.Note`.

## Errors

- None.

## Side Effects

- Only note (and update_time) should change; status remains success.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("SetNote: unexpected error: %v", err)
	}
	row := requireRow(t, resp)
	assertStatus(t, row, "success")
	if row.Note != req.Note {
		t.Fatalf("Note: got %q want %q", row.Note, req.Note)
	}
}
```
