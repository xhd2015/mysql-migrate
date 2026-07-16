# Scenario

**Feature**: `plan` with all-pending fixtures lists apply and exits 0

```
# clear chain — plan shows apply rows only
write two pending files, no logs
cli.Run(cfg, ["plan"]) -> apply for both, exit 0
```

## Preconditions

- Two non-EO fixtures; no log rows for those ids.
- Plan should list both as **apply** (non-skip).

## Steps

1. Write two fixtures; clear log ids.
2. Run plan.
3. Expect apply ×2, exit 0.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	const body = "SELECT 1;\n-- p5 plan all-pending\n"
	dir := t.TempDir()
	f1 := simpleFileName(1, fixtureSlug("plpend", "a"))
	f2 := simpleFileName(2, fixtureSlug("plpend", "b"))
	id1 := writeMigration(t, dir, f1, body)
	id2 := writeMigration(t, dir, f2, body)

	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	if err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	deleteLogIDs(t, db, id1, id2)

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id1, id2}
	req.Args = []string{"plan"}
	return nil
}
```
