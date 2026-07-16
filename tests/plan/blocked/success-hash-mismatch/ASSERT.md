## Expected

- `Run` succeeds.
- Item 0: Action **blocked**, HashMismatch=true, LogStatus `success`.
- Reason contains `hash_mismatch` (machine-readable).
- Item 1: Action **deferred**.
- `HasBlock` is true.

## Errors

- None.

```go
import "strings"

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
	it := resp.Items[0]
	if it.Action != "blocked" {
		t.Fatalf("Items[0].Action: got %q want blocked", it.Action)
	}
	if !it.HashMismatch {
		t.Fatal("Items[0].HashMismatch: got false want true")
	}
	if it.LogStatus != "success" {
		t.Fatalf("Items[0].LogStatus: got %q want success", it.LogStatus)
	}
	if !strings.Contains(strings.ToLower(it.Reason), "hash_mismatch") {
		t.Fatalf("Items[0].Reason: got %q; want substring hash_mismatch", it.Reason)
	}
	if resp.Items[1].Action != "deferred" {
		t.Fatalf("Items[1].Action: got %q want deferred", resp.Items[1].Action)
	}
}
```
