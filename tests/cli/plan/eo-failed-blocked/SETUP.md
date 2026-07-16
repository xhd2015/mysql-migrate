# Scenario

**Feature**: `plan` with EXACTLY-ONCE failed blocks chain and exits 1

```
# EO failed never auto re-applies; later pending deferred
files [EO-failed, later-pending]
seed failed(EO)
cli.Run(cfg, ["plan"])
  -> blocked EO, deferred later, exit 1
```

## Preconditions

- First file: `[EXACTLY-ONCE]` in filename; log status **failed** with matching hash.
- Second file: non-EO pending (no log).
- Plan output includes **blocked** and **deferred** (non-skip).

## Steps

1. Write EO + later fixtures; seed EO failed.
2. Run plan.
3. Expect exit 1, blocked then deferred on stdout.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	const (
		bodyEO = "SELECT 1;\n-- p5 plan eo-failed\n"
		bodyB  = "SELECT 2;\n-- p5 plan eo-failed later\n"
	)
	dir := t.TempDir()
	fEO := eoFileName(1, fixtureSlug("pleo", "drop"))
	fB := simpleFileName(2, fixtureSlug("pleo", "later"))
	idEO := writeMigration(t, dir, fEO, bodyEO)
	idB := writeMigration(t, dir, fB, bodyB)

	hashEO := contentSHA256(bodyEO)
	db := openLocalDB(t)
	t.Cleanup(func() { _ = db.Close() })
	seedFailed(t, db, idEO, true /* exactlyOnce */, hashEO, 7, "simulated EO failure")
	deleteLogIDs(t, db, idB)

	req.MigrationsDir = dir
	req.FixtureIDs = []string{idEO, idB}
	req.Args = []string{"plan"}
	return nil
}
```
