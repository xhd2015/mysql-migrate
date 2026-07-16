# Scenario

**Feature**: `apply --to <mid>` applies only through that migration_id inclusive

```
# three pending CREATE TABLEs; --to second id stops after second apply
files [A, B, C] all pending
cli.Run(cfg, ["apply", "--to", idB])
  -> apply A, apply B, stop; C remains pending / table C absent
```

## Preconditions

- Three valid non-EO CREATE TABLE fixtures (filename order A < B < C).
- No prior log rows for A/B/C.
- `Args` include `--to` equal to FixtureIDs[1] (mid).

## Steps

1. Write three CREATE TABLE files; drop tables; clear logs.
2. Run `apply --to <idB>`.
3. Expect A+B success + tables; C no success log and no table; exit 0.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	tblA := fixtureTable("apto", "a")
	tblB := fixtureTable("apto", "b")
	tblC := fixtureTable("apto", "c")
	bodyA := createTableSQL(tblA) + "-- p5 apply to-mid a\n"
	bodyB := createTableSQL(tblB) + "-- p5 apply to-mid b\n"
	bodyC := createTableSQL(tblC) + "-- p5 apply to-mid c\n"

	f1 := simpleFileName(1, fixtureSlug("apto", "a"))
	f2 := simpleFileName(2, fixtureSlug("apto", "b"))
	f3 := simpleFileName(3, fixtureSlug("apto", "c"))
	idA := writeMigration(t, dir, f1, bodyA)
	idB := writeMigration(t, dir, f2, bodyB)
	idC := writeMigration(t, dir, f3, bodyC)

	db := openLocalDB(t)
	t.Cleanup(func() {
		dropTables(t, db, tblA, tblB, tblC)
		deleteLogIDs(t, db, idA, idB, idC)
		_ = db.Close()
	})
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	dropTables(t, db, tblA, tblB, tblC)
	deleteLogIDs(t, db, idA, idB, idC)

	req.MigrationsDir = dir
	req.FixtureIDs = []string{idA, idB, idC}
	req.TableNames = []string{tblA, tblB, tblC}
	req.Args = []string{"apply", "--to", idB}
	return nil
}
```
