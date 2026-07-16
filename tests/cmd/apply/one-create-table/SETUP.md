# Scenario

**Feature**: binary `--dsn` + `--dir` apply one CREATE TABLE → success

```
# one pending migration; flags wire Config into cli.Run apply
write 01-<slug>.sql (CREATE TABLE IF NOT EXISTS)
mysql-migrate --dsn <harness> --dir <tmp> apply
  -> apply ok progress + log success + table exists; exit 0
```

## Preconditions

- MySQL reachable via harness DSN (else **skip**).
- One valid non-EO migration file with `CREATE TABLE IF NOT EXISTS`.
- Unique table name and migration_id for this session/leaf.
- No prior log row for that id; table dropped before run.

## Steps

1. Ensure MySQL; create temp dir; write one CREATE TABLE migration.
2. Delete prior log row / drop table.
3. Args: `--dsn`, harness DSN, `--dir`, dir, `apply`.
4. Expect exit 0, log success, table exists.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureMySQL(t)

	dir := t.TempDir()
	tbl := fixtureTable("oct", "a")
	body := createTableSQL(tbl) + "-- p6 cmd apply one-create\n"
	f1 := simpleFileName(1, fixtureSlug("oct", "a"))
	id1 := writeMigration(t, dir, f1, body)

	db := openLocalDB(t)
	t.Cleanup(func() {
		dropTables(t, db, tbl)
		deleteLogIDs(t, db, id1)
		_ = db.Close()
	})
	if err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	dropTables(t, db, tbl)
	deleteLogIDs(t, db, id1)

	dsn := harnessDSN()
	req.MigrationsDir = dir
	req.FixtureIDs = []string{id1}
	req.TableNames = []string{tbl}
	req.AssertDSN = dsn
	req.ClearMigrateEnv = true
	req.Args = []string{"--dsn", dsn, "--dir", dir, "apply"}
	return nil
}
```
