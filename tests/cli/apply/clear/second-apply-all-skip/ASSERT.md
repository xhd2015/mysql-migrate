## Expected

- Exit code **0**.
- Not stubbed.
- Both FixtureIDs remain log status **success**.
- Both TableNames still exist.
- No progress line claiming a new **failed** apply for either id.

## Side Effects

- No destructive re-exec required; tables/logs stable under skip.

## Errors

- stderr must not contain a hard Error about blocked/failed apply for these ids.

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
	requireNotStub(t, resp, "apply")

	if len(req.FixtureIDs) != 2 || len(req.TableNames) != 2 {
		t.Fatalf("setup: want 2 FixtureIDs and 2 TableNames, got ids=%v tables=%v",
			req.FixtureIDs, req.TableNames)
	}

	for _, id := range req.FixtureIDs {
		if stdoutHasApplyProgress(resp.Stdout, id, "failed") {
			t.Fatalf("second apply must not fail %q\nstdout:\n%s", id, resp.Stdout)
		}
	}
	if strings.Contains(strings.ToLower(resp.Stderr), "error:") &&
		strings.Contains(strings.ToLower(resp.Stderr), "blocked") {
		t.Fatalf("second apply should not refuse as blocked:\nstderr=%q", resp.Stderr)
	}

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	for _, id := range req.FixtureIDs {
		requireLogStatus(t, db, id, "success")
	}
	for _, tbl := range req.TableNames {
		if !tableExists(t, db, tbl) {
			t.Fatalf("expected table %q to still exist after second apply", tbl)
		}
	}
}
```
