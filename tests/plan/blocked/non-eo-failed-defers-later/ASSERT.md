## Expected

- `Run` succeeds.
- Item 0: Action **blocked**, ExactlyOnce=false, LogStatus `failed`.
- Item 1: Action **deferred**.
- `HasBlock` is true.

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
	if !resp.HasBlock {
		t.Fatal("HasBlock: got false want true")
	}
	if len(resp.Items) != 2 {
		t.Fatalf("len(Items): got %d want 2", len(resp.Items))
	}
	if resp.Items[0].Action != "blocked" {
		t.Fatalf("Items[0].Action: got %q want blocked", resp.Items[0].Action)
	}
	if resp.Items[0].ExactlyOnce {
		t.Fatal("Items[0].ExactlyOnce: got true want false")
	}
	if resp.Items[0].LogStatus != "failed" {
		t.Fatalf("Items[0].LogStatus: got %q want failed", resp.Items[0].LogStatus)
	}
	if resp.Items[1].Action != "deferred" {
		t.Fatalf("Items[1].Action: got %q want deferred", resp.Items[1].Action)
	}
}
```
