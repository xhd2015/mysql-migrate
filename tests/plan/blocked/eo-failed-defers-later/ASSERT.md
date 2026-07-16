## Expected

- `Run` succeeds.
- Item 0 (EO): Action **blocked**, ExactlyOnce=true, LogStatus `failed`.
- Item 1: Action **deferred** (would have been apply).
- `HasBlock` is true.
- Reason on item 0 indicates failure / exactly-once (contains `failed` or `exactly_once`).

## Errors

- None (plan always builds; block is data not error).

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
	eo, later := resp.Items[0], resp.Items[1]
	if !eo.ExactlyOnce {
		t.Fatal("Items[0].ExactlyOnce: got false want true")
	}
	if eo.Action != "blocked" {
		t.Fatalf("Items[0].Action: got %q want blocked", eo.Action)
	}
	if eo.LogStatus != "failed" {
		t.Fatalf("Items[0].LogStatus: got %q want failed", eo.LogStatus)
	}
	// Reason should be machine-readable around failed / exactly-once.
	r := strings.ToLower(eo.Reason)
	if !strings.Contains(r, "fail") && !strings.Contains(r, "exactly") {
		t.Fatalf("Items[0].Reason: got %q; want substring fail or exactly", eo.Reason)
	}
	if later.Action != "deferred" {
		t.Fatalf("Items[1].Action: got %q want deferred", later.Action)
	}
	if later.MigrationID != "2026-07-16-02-create-b" {
		t.Fatalf("Items[1].MigrationID: got %q", later.MigrationID)
	}
}
```
