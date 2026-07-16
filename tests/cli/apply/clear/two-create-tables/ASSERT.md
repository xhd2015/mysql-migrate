## Expected

- Exit code **0**.
- Not stubbed.
- Progress on stdout: each FixtureID co-occurs with **ok**.
- Summary mentions applied count (token `applied`).
- Log status **success** for both FixtureIDs.
- Both TableNames exist in the database.

## Side Effects

- `t_sql_migration_log` rows success for both ids.
- Both CREATE TABLE targets exist until Cleanup DROP.

## Errors

- stderr should not require Error/blocked (may be empty).

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
		if !stdoutHasApplyProgress(resp.Stdout, id, "ok") {
			t.Fatalf("want apply progress ok near %q\nstdout:\n%s", id, resp.Stdout)
		}
	}
	outLower := strings.ToLower(resp.Stdout)
	if !strings.Contains(outLower, "applied") {
		t.Fatalf("summary should mention applied:\n%s", resp.Stdout)
	}

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	for _, id := range req.FixtureIDs {
		requireLogStatus(t, db, id, "success")
	}
	for _, tbl := range req.TableNames {
		if !tableExists(t, db, tbl) {
			t.Fatalf("expected table %q to exist after apply", tbl)
		}
	}
}
```
