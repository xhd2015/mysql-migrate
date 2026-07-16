# Scenario

**Feature**: allow-retry on non-EO failed is a business error (exit 1)

```
# non-exactly-once must not use allow-retry clear path
seed failed(exactlyOnce=false)
cli.Run(cfg, ["allow-retry", id, "--note", "should fail"])
  -> exit 1; stderr Error; status remains failed
```

## Preconditions

- Seed failed with ExactlyOnce=false.
- Non-empty note (failure is EO gate, not note validation).

## Steps

1. Seed non-EO failed; run allow-retry.
2. Expect exit 1, Error on stderr, log still failed.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	const note = "p5 should not clear non-eo via allow-retry"
	dir := t.TempDir()
	body := "SELECT 1;\n-- p5 allow-retry non-eo\n"
	f := simpleFileName(1, fixtureSlug("arno", "ne"))
	id := writeMigration(t, dir, f, body)
	hash := contentSHA256(body)

	db := openLocalDB(t)
	t.Cleanup(func() {
		deleteLogIDs(t, db, id)
		_ = db.Close()
	})
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	seedFailed(t, db, id, false /* exactlyOnce */, hash, 12, "non-eo failed seed")

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id}
	req.RecoveryNote = note
	req.Args = []string{"allow-retry", id, "--note", note}
	return nil
}
```
