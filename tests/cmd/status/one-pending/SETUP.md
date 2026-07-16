# Scenario

**Feature**: binary `--dsn` + `--dir` status shows one pending migration

```
# one inventory file, no log row → plan action apply
write 01-<slug>.sql
mysql-migrate --dsn <harness> --dir <tmp> status
  -> status table with MIGRATION_ID + apply; exit 0
```

## Preconditions

- MySQL reachable via harness DSN (else **skip**).
- One valid non-EO migration file; no prior success log required.
- Unique migration_id for this session/leaf.
- Exclusive MySQL lock held (parent status Setup).

## Steps

1. Ensure MySQL + exclusive lock (parent); create temp dir; write one SQL file.
2. Args: `--dsn`, harness DSN, `--dir`, dir, `status`.
3. Expect exit 0, header row, fixture id with apply action.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	body := "-- p2 cmd status one-pending\nSELECT 1;\n"
	f1 := simpleFileName(1, fixtureSlug("stp", "a"))
	id1 := writeMigration(t, dir, f1, body)

	dsn := harnessDSN()
	req.MigrationsDir = dir
	req.FixtureIDs = []string{id1}
	req.AssertDSN = dsn
	req.ClearMigrateEnv = true
	req.Args = []string{"--dsn", dsn, "--dir", dir, "status"}
	return nil
}
```
