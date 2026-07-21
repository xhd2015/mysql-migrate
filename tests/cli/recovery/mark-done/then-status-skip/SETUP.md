# Scenario

**Feature**: mark-done with note forces success; follow-up status shows skip

```
# human verified migration offline — mark done, then status is skip (no re-apply)
seed failed(id, hash matches file)
cli.Run(cfg, ["mark-done", id, "--note", "ops verified"])
  -> exit 0; log success + note
FollowUp: status
  -> id action skip (hash match)
```

## Preconditions

- One non-EO fixture file; body stable so ContentSHA256 matches seed.
- Log seed: **failed** with that hash.
- `RecoveryNote` non-empty; Args include same note.
- FollowUpArgs: `status` with same Config (DSN + MigrationsDir).

## Steps

1. Write fixture; seed failed log; set mark-done Args + status FollowUp.
2. Expect primary exit 0, not stubbed, log success+note.
3. Expect follow-up status exit 0 and **skip** near id.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	fillConfigForDB(t, req)
	const note = "p5 ops verified offline — mark-done"
	dir := t.TempDir()
	body := "SELECT 1;\n-- p5 mark-done then-status-skip\n"
	f := simpleFileName(1, fixtureSlug(d, "mdone", "one"))
	id := writeMigration(t, dir, f, body)
	hash := contentSHA256(body)

	db := openLocalDB(t, d)
	t.Cleanup(func() {
		deleteLogIDs(t, db, id)
		_ = db.Close()
	})
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	seedFailed(t, db, id, false /* exactlyOnce */, hash, 15, "simulated fail before mark-done")

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id}
	req.RecoveryNote = note
	req.Args = []string{"mark-done", id, "--note", note}
	req.FollowUpArgs = []string{"status"}
	return nil
}
```
