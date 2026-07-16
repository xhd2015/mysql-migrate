## Expected

- `Run` succeeds.
- One item: Action **skip**, HashMismatch=false, LogStatus `success`.
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
	if len(resp.Items) != 1 {
		t.Fatalf("len(Items): got %d want 1", len(resp.Items))
	}
	it := resp.Items[0]
	if it.Action != "skip" {
		t.Fatalf("Action: got %q want skip", it.Action)
	}
	if it.HashMismatch {
		t.Fatal("HashMismatch: got true want false (empty log hash is not a mismatch)")
	}
	if it.LogStatus != "success" {
		t.Fatalf("LogStatus: got %q want success", it.LogStatus)
	}
}
```
