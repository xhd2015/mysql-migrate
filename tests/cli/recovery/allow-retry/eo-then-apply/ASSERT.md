## Expected

- Primary (`allow-retry`) exit **0**; not stubbed.
- Follow-up (`apply`) exit **0**; not stubbed in follow-up output.
- After full chain: log FixtureIDs[0] status **success**.
- TableNames[0] exists.
- Follow-up stdout has apply progress **ok** near the migration id.

## Side Effects

- EO failed cleared to pending, then applied once to success.

## Errors

- Primary and follow-up must not remain `not implemented`.

## Exit Code

- Primary 0; FollowUp 0

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
	requireNotStub(t, resp, "allow-retry")
	requireFollowUpExit(t, resp, 0)

	followCombined := strings.ToLower(resp.FollowUpStdout + "\n" + resp.FollowUpStderr)
	if strings.Contains(followCombined, "not implemented") {
		t.Fatalf("follow-up apply still stubbed:\nstdout=%q\nstderr=%q",
			resp.FollowUpStdout, resp.FollowUpStderr)
	}

	if len(req.FixtureIDs) != 1 || len(req.TableNames) != 1 {
		t.Fatalf("setup: want 1 FixtureID and 1 TableName, got ids=%v tables=%v",
			req.FixtureIDs, req.TableNames)
	}
	id := req.FixtureIDs[0]
	tbl := req.TableNames[0]

	if !stdoutHasApplyProgress(resp.FollowUpStdout, id, "ok") {
		t.Fatalf("follow-up apply must report ok for %q\nstdout:\n%s\nstderr:\n%s",
			id, resp.FollowUpStdout, resp.FollowUpStderr)
	}

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	requireLogStatus(t, db, id, "success")
	if !tableExists(t, db, tbl) {
		t.Fatalf("expected table %q after allow-retry then apply", tbl)
	}
}
```
