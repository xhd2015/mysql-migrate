## Expected

- Exit code **0**.
- FixtureIDs[0]: action **skip** on stdout (prior success, hash match).
- FixtureIDs[1]: action **apply** on stdout (no log).
- Seeded duration (`SuccessDurationMS`, e.g. `42`) appears on stdout.
- Not stubbed.

## Side Effects

- Seeded success row remains until Cleanup DELETE.

## Errors

- No `warning:` required (hash matches).

## Exit Code

- 0

```go
import (
	"fmt"
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
	idSkip, idApply := req.FixtureIDs[0], req.FixtureIDs[1]
	if !stdoutHasActionNearID(resp.Stdout, idSkip, "skip") {
		t.Fatalf("want skip near %q\nstdout:\n%s", idSkip, resp.Stdout)
	}
	if !stdoutHasActionNearID(resp.Stdout, idApply, "apply") {
		t.Fatalf("want apply near %q\nstdout:\n%s", idApply, resp.Stdout)
	}
	if req.SuccessDurationMS > 0 {
		wantDur := fmt.Sprintf("%d", req.SuccessDurationMS)
		if !strings.Contains(resp.Stdout, wantDur) {
			t.Fatalf("status should surface duration_ms %s on stdout:\n%s", wantDur, resp.Stdout)
		}
	}
}
```
