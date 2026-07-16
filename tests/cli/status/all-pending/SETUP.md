# Scenario

**Feature**: `status` with fixture files and no logs shows apply for each

```
# two pending migrations → both action apply, exit 0
write 01-a.sql, 02-b.sql (no log rows)
cli.Run(cfg, ["status"]) -> stdout apply for both ids, exit 0
```

## Preconditions

- Two valid non-EO migration files in temp dir (sorted by filename).
- Log table has **no** rows for those migration_ids (delete before run).
- Bodies are simple `SELECT 1;` (not executed by status).

## Steps

1. Write two fixtures; delete any prior log rows for those ids.
2. Run status.
3. Expect both ids with action **apply**, exit 0.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	const body = "SELECT 1;\n-- p5 status all-pending\n"
	dir := t.TempDir()
	f1 := simpleFileName(1, fixtureSlug("stpend", "a"))
	f2 := simpleFileName(2, fixtureSlug("stpend", "b"))
	id1 := writeMigration(t, dir, f1, body)
	id2 := writeMigration(t, dir, f2, body)

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	deleteLogIDs(t, db, id1, id2)

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id1, id2}
	req.Args = []string{"status"}
	return nil
}
```
