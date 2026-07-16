## Expected

- Exit code **1** (`HasBlock=true`).
- FixtureIDs[0]: action **blocked** on stdout.
- FixtureIDs[1]: action **deferred** on stdout (stop-chain after block).
- Stderr contains a line with **`warning:`** and the blocked migration id;
  preferably also mentions hash.
- Not stubbed.

## Side Effects

- Seeded success row until Cleanup; no SQL apply.

## Errors

- Business-gated plan (exit 1), not usage (2).

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
		t.Fatalf("status still stubbed:\nstdout=%q\nstderr=%q", resp.Stdout, resp.Stderr)
	}
	if len(req.FixtureIDs) != 2 {
		t.Fatalf("test setup: want 2 FixtureIDs, got %v", req.FixtureIDs)
	}
	idBlock, idDefer := req.FixtureIDs[0], req.FixtureIDs[1]
	if !stdoutHasActionNearID(resp.Stdout, idBlock, "blocked") {
		t.Fatalf("want blocked near %q\nstdout:\n%s", idBlock, resp.Stdout)
	}
	if !stdoutHasActionNearID(resp.Stdout, idDefer, "deferred") {
		t.Fatalf("want deferred near %q\nstdout:\n%s", idDefer, resp.Stdout)
	}

	if !strings.Contains(resp.Stderr, "warning:") {
		t.Fatalf("stderr must contain %q for hash mismatch:\nstderr=%q\nstdout=%q", "warning:", resp.Stderr, resp.Stdout)
	}
	if !strings.Contains(resp.Stderr, idBlock) {
		t.Fatalf("hash-mismatch warning should mention migration id %q:\nstderr=%q", idBlock, resp.Stderr)
	}
	if !strings.Contains(strings.ToLower(resp.Stderr), "hash") {
		t.Fatalf("hash-mismatch warning should mention hash:\nstderr=%q", resp.Stderr)
	}
}
```
