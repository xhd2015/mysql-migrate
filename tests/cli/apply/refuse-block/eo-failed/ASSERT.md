## Expected

- Exit code **1**.
- Not stubbed.
- stderr contains **Error** and **blocked**.
- FixtureIDs[0] (EO) remains **failed**.
- FixtureIDs[1] (later): not **success**.
- TableNames[0] does **not** exist.
- No apply **ok** progress for the later id.

## Side Effects

- No auto re-apply of EXACTLY-ONCE failed until allow-retry.
- Later deferred migration not executed.

## Errors

- Business refuse exit 1 (not usage 2).

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
	requireNotStub(t, resp, "apply")

	if len(req.FixtureIDs) != 2 || len(req.TableNames) != 1 {
		t.Fatalf("setup: want 2 FixtureIDs and 1 TableName, got ids=%v tables=%v",
			req.FixtureIDs, req.TableNames)
	}
	idEO, idLater := req.FixtureIDs[0], req.FixtureIDs[1]
	tblLater := req.TableNames[0]

	stderr := resp.Stderr
	if !strings.Contains(stderr, "Error") {
		t.Fatalf("refuse-block stderr must contain Error:\nstderr=%q\nstdout=%q", stderr, resp.Stdout)
	}
	if !strings.Contains(strings.ToLower(stderr), "blocked") {
		t.Fatalf("refuse-block stderr must mention blocked:\nstderr=%q", stderr)
	}

	if stdoutHasApplyProgress(resp.Stdout, idLater, "ok") {
		t.Fatalf("must not apply later %q when EO blocked\nstdout:\n%s", idLater, resp.Stdout)
	}

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	requireLogStatus(t, db, idEO, "failed")
	if st, ok := logStatus(t, db, idLater); ok && st == "success" {
		t.Fatalf("later %q must not be success when blocked (got %q)", idLater, st)
	}
	if tableExists(t, db, tblLater) {
		t.Fatalf("later table %q must not exist when apply refused for EO failed", tblLater)
	}
}
```
