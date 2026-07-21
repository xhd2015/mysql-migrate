## Expected

- Exit code **0**; not stubbed.
- Log FixtureIDs[0]: status still **success**, Note == RecoveryNote.

## Side Effects

- Note-only mutation; lifecycle status unchanged.

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
	requireNotStub(t, resp, "note")

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
}
```
