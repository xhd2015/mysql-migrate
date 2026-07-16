## Expected

- Exit code **0**.
- Stdout includes both fixture migration ids.
- Each id is associated with action **apply** (same line preferred).
- Not stubbed (`not implemented` must be absent).

## Side Effects

- EnsureTable may create log table; no MarkRunning from status itself.

## Errors

- stderr may be empty (no hash warnings).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	requireExit(t, resp, 0)

	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if strings.Contains(combined, "not implemented") {
		t.Fatalf("status still stubbed:\nstdout=%q\nstderr=%q", resp.Stdout, resp.Stderr)
	}
	if len(req.FixtureIDs) != 2 {
		t.Fatalf("test setup: want 2 FixtureIDs, got %v", req.FixtureIDs)
	}
	for _, id := range req.FixtureIDs {
		if !stdoutHasActionNearID(resp.Stdout, id, "apply") {
			t.Fatalf("status all-pending: want apply near %q\nstdout:\n%s", id, resp.Stdout)
		}
	}
}
```
