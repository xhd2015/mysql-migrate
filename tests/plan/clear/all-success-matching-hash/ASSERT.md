## Expected

- `Run` succeeds.
- Both items **skip**, `HasBlock=false`.
- `LogStatus` is `success` for each.
- `HashMismatch` is false.

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
	for i, it := range resp.Items {
		if it.Action != "skip" {
			t.Fatalf("Items[%d].Action: got %q want skip", i, it.Action)
		}
		if it.LogStatus != "success" {
			t.Fatalf("Items[%d].LogStatus: got %q want success", i, it.LogStatus)
		}
		if it.HashMismatch {
			t.Fatalf("Items[%d].HashMismatch: got true want false", i)
		}
	}
	if resp.Items[0].MigrationID != "2026-07-16-01-create-a" {
		t.Fatalf("Items[0].MigrationID: got %q", resp.Items[0].MigrationID)
	}
	if resp.Items[1].MigrationID != "2026-07-16-02-create-b" {
		t.Fatalf("Items[1].MigrationID: got %q", resp.Items[1].MigrationID)
	}
}
```
