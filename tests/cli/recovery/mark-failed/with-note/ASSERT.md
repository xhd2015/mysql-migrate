## Expected

- Exit code **0**; not stubbed.
- Log FixtureIDs[0]: status **failed**, Note == RecoveryNote.

## Side Effects

- Operator-forced failed with audit note; no migration SQL Exec required.

## Errors

- None on happy path.

## Exit Code

- 0

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	requireExit(t, resp, 0)
	requireNotStub(t, resp, "mark-failed")

	if len(req.FixtureIDs) != 1 {
		t.Fatalf("setup: want 1 FixtureID, got %v", req.FixtureIDs)
	}
	id := req.FixtureIDs[0]
	if req.RecoveryNote == "" {
		t.Fatal("setup: RecoveryNote empty")
	}

	db := openLocalDB(t, d)
	t.Cleanup(func() { _ = db.Close() })
	requireLogStatus(t, db, id, "failed")
	requireLogNote(t, db, id, req.RecoveryNote)
}
```
