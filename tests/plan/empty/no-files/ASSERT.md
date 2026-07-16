## Expected

- `Run` succeeds (`err == nil`).
- `Items` is empty (len 0).
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
	if len(resp.Items) != 0 {
		t.Fatalf("Items: got %d want 0 (%v)", len(resp.Items), resp.Items)
	}
	if resp.HasBlock {
		t.Fatal("HasBlock: got true want false")
	}
}
```
