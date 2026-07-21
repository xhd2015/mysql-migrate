# Scenario

**Feature**: note command updates note only; status remains success

```
# annotate without changing lifecycle status
seed success(id) with optional prior note
cli.Run(cfg, ["note", id, "--note", "updated ops note"])
  -> exit 0; Note=updated; status still success
```

## Preconditions

- Seed success; RecoveryNote is the new note string.
- Prior seed note differs from RecoveryNote so update is observable.

## Steps

1. Seed success with prior note; run note with new note.
2. Expect exit 0, status success, Note == RecoveryNote.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	const note = "p5 updated ops annotation"
	dir := t.TempDir()
	body := "SELECT 1;\n-- p5 note update-only\n"
	f := simpleFileName(1, fixtureSlug(d, "note", "one"))
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
	seedSuccess(t, db, id, false, hash, 5, "prior-note-before-update")

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id}
	req.RecoveryNote = note
	req.Args = []string{"note", id, "--note", note}
	return nil
}
```
