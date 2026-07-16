## Expected

- Exit code **1**.
- Not stubbed.
- FixtureIDs[0]: log status **failed**; progress contains **failed** near id.
- FixtureIDs[1]: **not** success.
- TableNames[0] does **not** exist.

## Side Effects

- MarkFailed on first; no successful apply of later migration.

## Errors

- Exit 1 (not usage 2).

## Exit Code

- 1

```go
import (
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
	idBad, idLater := req.FixtureIDs[0], req.FixtureIDs[1]
	tblLater := req.TableNames[0]

	if !stdoutHasApplyProgress(resp.Stdout, idBad, "failed") {
		t.Fatalf("want failed progress near %q\nstdout:\n%s\nstderr:\n%s",
			idBad, resp.Stdout, resp.Stderr)
	}
	if stdoutHasApplyProgress(resp.Stdout, idLater, "ok") {
		t.Fatalf("later migration %q must not report ok after stop\nstdout:\n%s", idLater, resp.Stdout)
	}

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	requireLogStatus(t, db, idBad, "failed")
	if st, ok := logStatus(t, db, idLater); ok && st == "success" {
		t.Fatalf("later %q must not be success after stop (got %q)", idLater, st)
	}
	if tableExists(t, db, tblLater) {
		t.Fatalf("later table %q must not exist when apply stopped on bad SQL", tblLater)
	}
}
```
