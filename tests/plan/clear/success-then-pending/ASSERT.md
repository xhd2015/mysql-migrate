## Expected

- `Run` succeeds.
- Item 0: MigrationID A, Action **skip**, LogStatus `success`.
- Item 1: MigrationID B, Action **apply**.
- `HasBlock` is false.

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
	a, b := resp.Items[0], resp.Items[1]
	if a.MigrationID != "2026-07-16-01-create-a" || a.Action != "skip" {
		t.Fatalf("item0: id=%q action=%q want id=...-a action=skip", a.MigrationID, a.Action)
	}
	if a.LogStatus != "success" {
		t.Fatalf("item0.LogStatus: got %q want success", a.LogStatus)
	}
	if b.MigrationID != "2026-07-16-02-create-b" || b.Action != "apply" {
		t.Fatalf("item1: id=%q action=%q want id=...-b action=apply", b.MigrationID, b.Action)
	}
}
```
