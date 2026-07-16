## Expected

- `Run` succeeds.
- Both items Action **apply**.
- `HasBlock` is false.
- Effective LogStatus reflects input (`""` / `"pending"`) — not rewritten to blocked.

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
	if resp.Items[0].Action != "apply" {
		t.Fatalf("Items[0].Action: got %q want apply (empty status)", resp.Items[0].Action)
	}
	if resp.Items[1].Action != "apply" {
		t.Fatalf("Items[1].Action: got %q want apply (pending status)", resp.Items[1].Action)
	}
	if resp.Items[0].LogStatus != "" {
		t.Fatalf("Items[0].LogStatus: got %q want empty", resp.Items[0].LogStatus)
	}
	if resp.Items[1].LogStatus != "pending" {
		t.Fatalf("Items[1].LogStatus: got %q want pending", resp.Items[1].LogStatus)
	}
}
```
