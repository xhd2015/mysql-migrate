# Scenario

**Feature**: mark-failed with note sets failed status and operator note in log

```
# force failed for audit without re-running SQL
seed success(id)
cli.Run(cfg, ["mark-failed", id, "--note", "ops abort"])
  -> exit 0; log status=failed, note set
```

## Preconditions

- Seed **success** row with matching hash for a fixture file.
- RecoveryNote non-empty.

## Steps

1. Seed success; run mark-failed --note.
2. Expect exit 0, not stubbed, log failed + note.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	const note = "p5 ops abort — mark-failed"
	dir := t.TempDir()
	body := "SELECT 1;\n-- p5 mark-failed with-note\n"
	f := simpleFileName(1, fixtureSlug(d, "mfail", "one"))
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
	seedSuccess(t, db, id, false, hash, 8, "prior success before mark-failed")

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id}
	req.RecoveryNote = note
	req.Args = []string{"mark-failed", id, "--note", note}
	return nil
}
```
