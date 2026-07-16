## Expected

- Exit code **0**.
- Not stubbed.
- FixtureIDs[0] and [1]: log **success**; TableNames[0] and [1] exist.
- FixtureIDs[2]: **not** success; TableNames[2] does **not** exist.
- Progress ok for first two ids; no ok progress for third.
- Summary should mention `pending`.

## Side Effects

- Only migrations up to and including `--to` id are applied this run.

## Errors

- No blocked Error required on stderr.

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

	if len(req.FixtureIDs) != 3 || len(req.TableNames) != 3 {
		t.Fatalf("setup: want 3 FixtureIDs and 3 TableNames, got ids=%v tables=%v",
			req.FixtureIDs, req.TableNames)
	}
	idA, idB, idC := req.FixtureIDs[0], req.FixtureIDs[1], req.FixtureIDs[2]
	tblA, tblB, tblC := req.TableNames[0], req.TableNames[1], req.TableNames[2]

	for _, id := range []string{idA, idB} {
		if !stdoutHasApplyProgress(resp.Stdout, id, "ok") {
			t.Fatalf("want apply ok near %q\nstdout:\n%s", id, resp.Stdout)
		}
	}
	if stdoutHasApplyProgress(resp.Stdout, idC, "ok") {
		t.Fatalf("--to mid must not apply %q\nstdout:\n%s", idC, resp.Stdout)
	}

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	requireLogStatus(t, db, idA, "success")
	requireLogStatus(t, db, idB, "success")
	if st, ok := logStatus(t, db, idC); ok && st == "success" {
		t.Fatalf("--to mid: %q must not be success (got %q)", idC, st)
	}
	if !tableExists(t, db, tblA) || !tableExists(t, db, tblB) {
		t.Fatalf("expected tables %q and %q to exist", tblA, tblB)
	}
	if tableExists(t, db, tblC) {
		t.Fatalf("table %q for migration after --to must not exist", tblC)
	}
	if !strings.Contains(strings.ToLower(resp.Stdout), "pending") {
		t.Fatalf("summary should mention pending for remaining migrations:\n%s", resp.Stdout)
	}
}
```
