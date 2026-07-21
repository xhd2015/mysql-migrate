# Scenario

**Feature**: bad SQL fails first migration; later pending is not applied

```
# first file invalid SQL; second CREATE TABLE pending
cli.Run(cfg, ["apply"])
  -> MarkRunning(id1) -> Exec fail -> MarkFailed
  -> STOP; id2 not success; table2 absent; exit 1
```

## Preconditions

- File 01: intentionally invalid SQL.
- File 02: valid `CREATE TABLE IF NOT EXISTS`.
- No prior log rows for either id; later table dropped before run.

## Steps

1. Write bad + good fixtures; clear logs/tables.
2. Run apply.
3. Expect exit 1, id1 failed, id2 not success, later table missing.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := t.TempDir()
	tblLater := fixtureTable(d, "apbad", "later")
	bodyBad := "THIS IS NOT VALID MYSQL SQL FOR P5 APPLY FAIL;\n"
	bodyOK := createTableSQL(tblLater) + "-- p5 apply bad-sql later\n"

	f1 := simpleFileName(1, fixtureSlug(d, "apbad", "bad"))
	f2 := simpleFileName(2, fixtureSlug(d, "apbad", "later"))
	idBad := writeMigration(t, dir, f1, bodyBad)
	idLater := writeMigration(t, dir, f2, bodyOK)

	db := openLocalDB(t, d)
	t.Cleanup(func() {
		dropTables(t, db, tblLater)
		deleteLogIDs(t, db, idBad, idLater)
		_ = db.Close()
	})
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	dropTables(t, db, tblLater)
	deleteLogIDs(t, db, idBad, idLater)

	req.MigrationsDir = dir
	req.FixtureIDs = []string{idBad, idLater}
	req.TableNames = []string{tblLater}
	req.Args = []string{"apply"}
	return nil
}
```
