# Scenario

**Feature**: second apply after success is all-skip, exit 0, no error

```
# prior success matching hash → plan skip for all; apply re-run is no-op exit 0
seed success(idA), success(idB); tables already created
cli.Run(cfg, ["apply"])
  -> no Action==apply items → exit 0; logs stay success
```

## Preconditions

- Two fixtures with known bodies; log **success** + matching ContentSHA256 for both.
- Tables pre-created (post-first-apply world).

## Steps

1. Write two CREATE TABLE files; create tables; seed success for both ids.
2. Run `apply` again.
3. Expect exit 0, not stubbed, both still success, tables still exist.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	tblA := fixtureTable("ap2nd", "a")
	tblB := fixtureTable("ap2nd", "b")
	bodyA := createTableSQL(tblA) + "-- p5 apply second-skip a\n"
	bodyB := createTableSQL(tblB) + "-- p5 apply second-skip b\n"

	f1 := simpleFileName(1, fixtureSlug("ap2nd", "a"))
	f2 := simpleFileName(2, fixtureSlug("ap2nd", "b"))
	id1 := writeMigration(t, dir, f1, bodyA)
	id2 := writeMigration(t, dir, f2, bodyB)

	hashA := contentSHA256(bodyA)
	hashB := contentSHA256(bodyB)

	db := openLocalDB(t)
	t.Cleanup(func() {
		dropTables(t, db, tblA, tblB)
		deleteLogIDs(t, db, id1, id2)
		_ = db.Close()
	})
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	dropTables(t, db, tblA, tblB)
	if _, err := db.Exec(createTableSQL(tblA)); err != nil {
		t.Fatalf("pre-create %s: %v", tblA, err)
	}
	if _, err := db.Exec(createTableSQL(tblB)); err != nil {
		t.Fatalf("pre-create %s: %v", tblB, err)
	}
	seedSuccess(t, db, id1, false, hashA, 11, "first-apply-ok")
	seedSuccess(t, db, id2, false, hashB, 12, "first-apply-ok")

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id1, id2}
	req.TableNames = []string{tblA, tblB}
	req.Args = []string{"apply"}
	return nil
}
```
