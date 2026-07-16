# Scenario

**Feature**: EO failed → allow-retry + note → apply runs once → success

```
# exit-criteria chain
files [EO CREATE TABLE]
seed failed(EO, hash matches)
cli.Run(cfg, ["allow-retry", id, "--note", "ops approved retry"])
  -> exit 0; status pending + note
FollowUp: apply
  -> MarkRunning -> Exec CREATE -> MarkSuccess; exit 0; table exists
```

## Preconditions

- Filename includes `[EXACTLY-ONCE]`; log seed failed + ExactlyOnce=true + matching hash.
- SQL body is idempotent CREATE TABLE so apply succeeds after clear.
- Table dropped before run; RecoveryNote non-empty.
- FollowUpArgs: `apply`.

## Steps

1. Write EO CREATE fixture; seed failed EO; drop table.
2. Primary allow-retry; follow-up apply.
3. Expect primary 0 + not stubbed; after chain log success + table exists + apply ok.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	const note = "p5 ops approved EO retry after fix"
	dir := t.TempDir()
	tbl := fixtureTable("areo", "t")
	body := createTableSQL(tbl) + "-- p5 allow-retry eo-then-apply\n"
	f := eoFileName(1, fixtureSlug("areo", "once"))
	id := writeMigration(t, dir, f, body)
	hash := contentSHA256(body)

	db := openLocalDB(t)
	t.Cleanup(func() {
		dropTables(t, db, tbl)
		deleteLogIDs(t, db, id)
		_ = db.Close()
	})
	if err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	dropTables(t, db, tbl)
	seedFailed(t, db, id, true /* exactlyOnce */, hash, 21, "simulated EO failure before allow-retry")

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id}
	req.TableNames = []string{tbl}
	req.RecoveryNote = note
	req.Args = []string{"allow-retry", id, "--note", note}
	req.FollowUpArgs = []string{"apply"}
	return nil
}
```
