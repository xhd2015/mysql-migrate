## Expected

- `err == nil` (no duplicate-key failure on second MarkRunning).
- `UniqueCount == 1` for this migration_id.
- Row found with second-call fields:
  - ExactlyOnce=true
  - ContentSHA256=SecondContentSHA256
  - AppliedBy=SecondAppliedBy
  - Status=`running` (MarkRunning leaves row running)

## Errors

- Must not surface MySQL duplicate entry / unique constraint errors.

## Side Effects

- Still exactly one row for migration_id after two MarkRunning calls.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unique_upsert: second MarkRunning must not error (prefer UPSERT): %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.UniqueCount != 1 {
		t.Fatalf("UniqueCount: got %d want 1 (no duplicate rows)", resp.UniqueCount)
	}
	row := requireRow(t, resp)
	assertStatus(t, row, "running")
	if !row.ExactlyOnce {
		t.Fatal("ExactlyOnce: got false want true (second MarkRunning)")
	}
	if row.ContentSHA256 != req.SecondContentSHA256 {
		t.Fatalf("ContentSHA256: got %q want %q", row.ContentSHA256, req.SecondContentSHA256)
	}
	if row.AppliedBy != req.SecondAppliedBy {
		t.Fatalf("AppliedBy: got %q want %q", row.AppliedBy, req.SecondAppliedBy)
	}
}
```
