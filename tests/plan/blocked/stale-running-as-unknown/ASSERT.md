## Expected

- `Run` succeeds.
- Items 0 and 1: Action **blocked**, effective LogStatus **`unknown`** (normalized from `running`).
- Item 0 ExactlyOnce=false; item 1 ExactlyOnce=true.
- Item 2: Action **deferred**.
- `HasBlock` is true.
- Reason may mention `stale` / `running` / `unknown`.

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
	if len(resp.Items) != 3 {
		t.Fatalf("len(Items): got %d want 3", len(resp.Items))
	}
	for i := 0; i < 2; i++ {
		it := resp.Items[i]
		if it.Action != "blocked" {
			t.Fatalf("Items[%d].Action: got %q want blocked", i, it.Action)
		}
		if it.LogStatus != "unknown" {
			t.Fatalf("Items[%d].LogStatus: got %q want unknown (stale running normalized)", i, it.LogStatus)
		}
	}
	if resp.Items[0].ExactlyOnce {
		t.Fatal("Items[0].ExactlyOnce: want false")
	}
	if !resp.Items[1].ExactlyOnce {
		t.Fatal("Items[1].ExactlyOnce: want true")
	}
	if resp.Items[2].Action != "deferred" {
		t.Fatalf("Items[2].Action: got %q want deferred", resp.Items[2].Action)
	}
	// At least one blocked reason should acknowledge stale/running/unknown.
	r0 := strings.ToLower(resp.Items[0].Reason)
	if !strings.Contains(r0, "unknown") && !strings.Contains(r0, "running") && !strings.Contains(r0, "stale") {
		t.Fatalf("Items[0].Reason: got %q; want unknown|running|stale", resp.Items[0].Reason)
	}
}
```
