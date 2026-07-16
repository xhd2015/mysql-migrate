## Expected

- Exit code **0**.
- Progress on stdout: FixtureID co-occurs with **ok** (or summary with `applied`).
- Log status **success** for FixtureIDs[0].
- TableNames[0] exists in the database.

## Side Effects

- `t_sql_migration_log` row success for the migration id.
- CREATE TABLE target exists until Cleanup DROP.

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

	if len(req.FixtureIDs) != 1 || len(req.TableNames) != 1 {
		t.Fatalf("setup: want 1 FixtureID and 1 TableName, got ids=%v tables=%v",
			req.FixtureIDs, req.TableNames)
	}
	id := req.FixtureIDs[0]
	tbl := req.TableNames[0]

	if !stdoutHasApplyProgress(resp.Stdout, id, "ok") {
		// fallback: summary still proves apply path ran
		if !strings.Contains(strings.ToLower(resp.Stdout), "applied") {
			t.Fatalf("want apply progress ok near %q or applied summary\nstdout:\n%s\nstderr:\n%s",
				id, resp.Stdout, resp.Stderr)
		}
	}

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	requireLogStatus(t, db, id, "success")
	if !tableExists(t, db, tbl) {
		t.Fatalf("expected table %q to exist after apply", tbl)
	}
}
```
