## Expected

- Exit code **1** (`HasBlock=true`).
- FixtureIDs[0] (EO): action **blocked** on stdout.
- FixtureIDs[1] (later): action **deferred** on stdout.
- Not stubbed.

## Side Effects

- Seeded failed EO row until Cleanup; no SQL apply.

## Errors

- Business failure exit 1 (not usage 2).

## Exit Code

- 1

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	requireExit(t, resp, 1)

	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if strings.Contains(combined, "not implemented") {
		t.Fatalf("plan still stubbed:\nstdout=%q\nstderr=%q", resp.Stdout, resp.Stderr)
	}
	if len(req.FixtureIDs) != 2 {
		t.Fatalf("test setup: want 2 FixtureIDs, got %v", req.FixtureIDs)
	}
	idEO, idLater := req.FixtureIDs[0], req.FixtureIDs[1]
	if !stdoutHasActionNearID(resp.Stdout, idEO, "blocked") {
		t.Fatalf("want blocked near EO id %q\nstdout:\n%s", idEO, resp.Stdout)
	}
	if !stdoutHasActionNearID(resp.Stdout, idLater, "deferred") {
		t.Fatalf("want deferred near %q\nstdout:\n%s", idLater, resp.Stdout)
	}
}
```
