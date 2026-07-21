# Scenario

**Feature**: success log with content hash ≠ file hash → blocked, exit 1, stderr warning

```
# file changed after success — never silent re-apply
seed success(idA, hash=old); file hash=new
cli.Run(cfg, ["status"])
  -> blocked idA, optional deferred later, exit 1
  -> stderr line "warning:" + migration id
```

## Preconditions

- Two fixtures: idA (success log with **wrong** hash), idB pending.
- File body for idA differs from seeded ContentSHA256.
- Status shows all items: blocked then deferred (stop-chain).

## Steps

1. Write files; seed success for idA with `deadbeef…` hash (not file hash).
2. Run status.
3. Expect exit 1, blocked on idA, deferred on idB, stderr `warning:`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	const (
		bodyA     = "SELECT 1;\n-- p5 status hash-mismatch current file\n"
		bodyB     = "SELECT 2;\n-- p5 status hash-mismatch later\n"
		wrongHash = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
	)
	dir := t.TempDir()
	f1 := simpleFileName(1, fixtureSlug(d, "sthash", "a"))
	f2 := simpleFileName(2, fixtureSlug(d, "sthash", "b"))
	id1 := writeMigration(t, dir, f1, bodyA)
	id2 := writeMigration(t, dir, f2, bodyB)

	db := openLocalDB(t, d)
	t.Cleanup(func() { _ = db.Close() })
	seedSuccess(t, db, id1, false, wrongHash, 10, "")
	deleteLogIDs(t, db, id2)

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id1, id2}
	req.Args = []string{"status"}
	return nil
}
```
