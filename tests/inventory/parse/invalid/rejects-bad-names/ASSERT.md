## Expected

- `Run` returns `err == nil` (batch completed).
- `ParseErrors` has one entry per `InvalidNames` element.
- Every `ParseErrors[i]` is non-empty (each name failed ParseFileName).

## Errors

- Each invalid basename surfaces a non-nil parse error (message text not pinned).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("batch parse: unexpected Run error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if len(resp.ParseErrors) != len(req.InvalidNames) {
		t.Fatalf("ParseErrors len: got %d want %d", len(resp.ParseErrors), len(req.InvalidNames))
	}
	for i, name := range req.InvalidNames {
		if resp.ParseErrors[i] == "" {
			t.Fatalf("name %q: expected parse error, got success", name)
		}
	}
}
```
