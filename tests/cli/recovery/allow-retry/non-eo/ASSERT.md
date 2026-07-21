## Expected

- Exit code **1** (biz error, not usage 2).
- Not stubbed.
- stderr contains **Error**.
- Log FixtureIDs[0] remains **failed**.

## Side Effects

- Row not cleared to pending; allow-retry rejected.

## Errors

- Business reject of non-EO allow-retry.

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
	requireNotStub(t, resp, "allow-retry")

	if !strings.Contains(resp.Stderr, "Error") {
		t.Fatalf("allow-retry non-EO stderr must contain Error:\nstderr=%q\nstdout=%q",
			resp.Stderr, resp.Stdout)
	}

	if len(req.FixtureIDs) != 1 {
		t.Fatalf("setup: want 1 FixtureID, got %v", req.FixtureIDs)
	}
	id := req.FixtureIDs[0]

	db := openLocalDB(t, d)
	t.Cleanup(func() { _ = db.Close() })
	requireLogStatus(t, db, id, "failed")
}
```
