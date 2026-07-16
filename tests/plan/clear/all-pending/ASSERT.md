## Expected

- `Run` succeeds.
- Two plan items in file order:
  1. `2026-07-16-01-create-a` → **apply**, ExactlyOnce=false
  2. `2026-07-16-02-create-b` → **apply**, ExactlyOnce=false
- `HasBlock` is false.
- Neither item has `HashMismatch`.

## Errors

- None.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("plan.Build: unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.HasBlock {
		t.Fatal("HasBlock: got true want false")
	}
	if len(resp.Items) != 2 {
		t.Fatalf("len(Items): got %d want 2", len(resp.Items))
	}
	wantIDs := []string{"2026-07-16-01-create-a", "2026-07-16-02-create-b"}
	for i, id := range wantIDs {
		it := resp.Items[i]
		if it.MigrationID != id {
			t.Fatalf("Items[%d].MigrationID: got %q want %q", i, it.MigrationID, id)
		}
		if it.Action != "apply" {
			t.Fatalf("Items[%d].Action: got %q want apply", i, it.Action)
		}
		if it.ExactlyOnce {
			t.Fatalf("Items[%d].ExactlyOnce: got true want false", i)
		}
		if it.HashMismatch {
			t.Fatalf("Items[%d].HashMismatch: got true want false", i)
		}
	}
}
```
