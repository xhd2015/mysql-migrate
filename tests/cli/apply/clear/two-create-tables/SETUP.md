# Scenario

**Feature**: apply two idempotent CREATE TABLE migrations → success + tables exist

```
# two pending CREATE TABLE IF NOT EXISTS → both applied, log success, exit 0
write 01-a.sql, 02-b.sql (no log rows)
cli.Run(cfg, ["apply"])
  -> MarkRunning+Exec+MarkSuccess for each
  -> tables exist; stdout ok progress + summary; exit 0
```

## Preconditions

- Two valid non-EO migration files with `CREATE TABLE IF NOT EXISTS` bodies.
- Unique table names per session (TableNames[0], TableNames[1]).
- Log table has **no** prior rows for those migration_ids.
- Tables dropped before run so existence proves Exec ran.

## Steps

1. Ensure MySQL; drop leftover fixture tables; write two CREATE TABLE files.
2. Delete prior log rows for those ids.
3. Run `apply`.
4. Expect both logs **success**, both tables exist, exit 0, progress ok.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	tblA := fixtureTable("ap2ct", "a")
	tblB := fixtureTable("ap2ct", "b")
	bodyA := createTableSQL(tblA) + "-- p5 apply two-create a\n"
	bodyB := createTableSQL(tblB) + "-- p5 apply two-create b\n"

	f1 := simpleFileName(1, fixtureSlug("ap2ct", "a"))
	f2 := simpleFileName(2, fixtureSlug("ap2ct", "b"))
	id1 := writeMigration(t, dir, f1, bodyA)
	id2 := writeMigration(t, dir, f2, bodyB)

	db := openLocalDB(t)
	t.Cleanup(func() {
		dropTables(t, db, tblA, tblB)
		deleteLogIDs(t, db, id1, id2)
		_ = db.Close()
	})
	if err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	dropTables(t, db, tblA, tblB)
	deleteLogIDs(t, db, id1, id2)

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id1, id2}
	req.TableNames = []string{tblA, tblB}
	req.Args = []string{"apply"}
	return nil
}
```
