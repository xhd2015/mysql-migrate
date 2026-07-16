## Expected

- `err == nil`, row found.
- `status == "success"`.
- `DurationMS == req.DurationMS` (100).
- `ContentSHA256 == req.ContentSHA256`.
- `AppliedBy == req.AppliedBy`.
- `ExactlyOnce == false`.
- `MigrationID == req.MigrationID`.

## Errors

- None.

## Side Effects

- One row in `t_sql_migration_log` for this migration_id (success).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("lifecycle_success: unexpected error: %v", err)
	}
	row := requireRow(t, resp)
	if row.MigrationID != req.MigrationID {
		t.Fatalf("MigrationID: got %q want %q", row.MigrationID, req.MigrationID)
	}
	assertStatus(t, row, "success")
	if row.DurationMS != req.DurationMS {
		t.Fatalf("DurationMS: got %d want %d", row.DurationMS, req.DurationMS)
	}
	if row.ContentSHA256 != req.ContentSHA256 {
		t.Fatalf("ContentSHA256: got %q want %q", row.ContentSHA256, req.ContentSHA256)
	}
	if row.AppliedBy != req.AppliedBy {
		t.Fatalf("AppliedBy: got %q want %q", row.AppliedBy, req.AppliedBy)
	}
	if row.ExactlyOnce {
		t.Fatal("ExactlyOnce: got true want false")
	}
}
```
