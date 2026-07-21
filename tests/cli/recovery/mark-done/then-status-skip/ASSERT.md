## Expected

- Primary exit code **0**; not stubbed.
- Log for FixtureIDs[0]: status **success**, Note == RecoveryNote.
- Follow-up **status** ran with exit **0**.
- Follow-up stdout shows **skip** near the migration id.

## Side Effects

- Log forced to success with operator note for audit.
- Status path treats matching-hash success as skip (no apply).

## Errors

- stderr must not leave the command stubbed.

## Exit Code

- Primary 0; FollowUp 0

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	requireExit(t, resp, 0)
	requireNotStub(t, resp, "mark-done")
	requireFollowUpExit(t, resp, 0)

	if len(req.FixtureIDs) != 1 {
		t.Fatalf("setup: want 1 FixtureID, got %v", req.FixtureIDs)
	}
	id := req.FixtureIDs[0]
	if req.RecoveryNote == "" {
		t.Fatal("setup: RecoveryNote empty")
	}

	db := openLocalDB(t, d)
	t.Cleanup(func() { _ = db.Close() })
	requireLogStatus(t, db, id, "success")
	requireLogNote(t, db, id, req.RecoveryNote)

	if !stdoutHasActionNearID(resp.FollowUpStdout, id, "skip") {
		t.Fatalf("follow-up status must show skip near %q\nstdout:\n%s\nstderr:\n%s",
			id, resp.FollowUpStdout, resp.FollowUpStderr)
	}
}
```
