## Expected

- `err == nil`.
- `resp.Rows` contains every seeded migration_id (primary + extras).
- Each seeded row has status `running`.

## Errors

- None.

## Side Effects

- Three running rows left for session-prefixed ids (isolated from other tests).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("list: unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	wantIDs := append([]string{req.MigrationID}, req.ExtraMigrationIDs...)
	byID := make(map[string]RowView, len(resp.Rows))
	for _, r := range resp.Rows {
		byID[r.MigrationID] = r
	}
	for _, id := range wantIDs {
		row, ok := byID[id]
		if !ok {
			t.Fatalf("List missing migration_id %q (got %d rows total)", id, len(resp.Rows))
		}
		if row.Status != "running" {
			t.Fatalf("id %q status: got %q want running", id, row.Status)
		}
	}
}
```
