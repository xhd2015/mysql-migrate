# Scenario

**Feature**: `status` shows skip for success+matching hash and apply for pending

```
# first file success log (hash match) → skip; second no log → apply; exit 0
seed MarkRunning+MarkSuccess(idA, hash=file)
cli.Run(cfg, ["status"]) -> skip idA, apply idB, exit 0
```

## Preconditions

- Two fixtures: idA then idB (filename order).
- idA: log status success with **matching** ContentSHA256 and known duration_ms.
- idB: no log row.
- Status prints **all** items (skip + apply).

## Steps

1. Write two files with known bodies; seed success for idA only.
2. Run status.
3. Expect skip for idA (duration visible), apply for idB, exit 0.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	const (
		bodyA = "SELECT 1;\n-- p5 status success-a\n"
		bodyB = "SELECT 2;\n-- p5 status pending-b\n"
	)
	dir := t.TempDir()
	f1 := simpleFileName(1, fixtureSlug(d, "stskip", "a"))
	f2 := simpleFileName(2, fixtureSlug(d, "stskip", "b"))
	id1 := writeMigration(t, dir, f1, bodyA)
	id2 := writeMigration(t, dir, f2, bodyB)

	hashA := contentSHA256(bodyA)
	const durationMS = 42
	db := openLocalDB(t, d)
	t.Cleanup(func() { _ = db.Close() })
	seedSuccess(t, db, id1, false, hashA, durationMS, "seeded-ok")
	deleteLogIDs(t, db, id2) // ensure pending

	req.MigrationsDir = dir
	req.FixtureIDs = []string{id1, id2}
	req.SuccessDurationMS = durationMS
	req.Args = []string{"status"}
	return nil
}
```
