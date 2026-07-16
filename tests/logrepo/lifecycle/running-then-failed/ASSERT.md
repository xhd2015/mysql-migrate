## Expected

- `err == nil`, row found.
- `status == "failed"`.
- `ErrorMessage == req.ErrorMessage`.
- `DurationMS == req.DurationMS`.
- Hash and AppliedBy preserved from MarkRunning.

## Errors

- None.

## Side Effects

- One failed row for this migration_id.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("lifecycle_failed: unexpected error: %v", err)
	}
	row := requireRow(t, resp)
	assertStatus(t, row, "failed")
	if row.ErrorMessage != req.ErrorMessage {
		t.Fatalf("ErrorMessage: got %q want %q", row.ErrorMessage, req.ErrorMessage)
	}
	if row.DurationMS != req.DurationMS {
		t.Fatalf("DurationMS: got %d want %d", row.DurationMS, req.DurationMS)
	}
	if row.ContentSHA256 != req.ContentSHA256 {
		t.Fatalf("ContentSHA256: got %q want %q", row.ContentSHA256, req.ContentSHA256)
	}
	if row.AppliedBy != req.AppliedBy {
		t.Fatalf("AppliedBy: got %q want %q", row.AppliedBy, req.AppliedBy)
	}
}
```
