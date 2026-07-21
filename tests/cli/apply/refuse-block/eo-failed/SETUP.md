# Scenario

**Feature**: EXACTLY-ONCE failed blocks apply; later pending is not applied

```
# EO failed never auto re-applies; later CREATE must not run until allow-retry
files [EO-failed, later-pending]
seed failed(EO, exactlyOnce=true)
cli.Run(cfg, ["apply"])
  -> refuse; stderr Error blocked; later table absent; exit 1
```

## Preconditions

- First file: `[EXACTLY-ONCE]` in filename; log status **failed** with matching hash.
- Second file: non-EO pending CREATE TABLE; table dropped before run.

## Steps

1. Write EO + later fixtures; seed EO failed; clear later log/table.
2. Run apply.
3. Expect exit 1, stderr Error + blocked, later not success / no table.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	const bodyEO = "SELECT 1;\n-- p5 apply eo-failed block\n"
	dir := t.TempDir()
	tblLater := fixtureTable(d, "apeo", "later")
	bodyLater := createTableSQL(tblLater) + "-- p5 apply eo-failed later\n"

	fEO := eoFileName(1, fixtureSlug(d, "apeo", "drop"))
	fLater := simpleFileName(2, fixtureSlug(d, "apeo", "later"))
	idEO := writeMigration(t, dir, fEO, bodyEO)
	idLater := writeMigration(t, dir, fLater, bodyLater)

	hashEO := contentSHA256(bodyEO)
	db := openLocalDB(t, d)
	t.Cleanup(func() {
		dropTables(t, db, tblLater)
		deleteLogIDs(t, db, idEO, idLater)
		_ = db.Close()
	})
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	dropTables(t, db, tblLater)
	seedFailed(t, db, idEO, true /* exactlyOnce */, hashEO, 9, "simulated EO failure for apply")
	deleteLogIDs(t, db, idLater)

	req.MigrationsDir = dir
	req.FixtureIDs = []string{idEO, idLater}
	req.TableNames = []string{tblLater}
	req.Args = []string{"apply"}
	return nil
}
```
