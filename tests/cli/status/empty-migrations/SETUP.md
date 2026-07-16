# Scenario

**Feature**: `status` with empty migrations dir exits 0 (nothing to do)

```
# empty inventory → empty plan, HasBlock=false
cfg.MigrationsDir=<empty tmp>
cli.Run(cfg, ["status"]) -> exit 0, no apply/blocked/deferred data
```

## Preconditions

- Temp migrations dir exists and has **zero** grammar-matching `.sql` files.
- No log seeding required.
- MySQL still opened / EnsureTable still runs.

## Steps

1. Create empty temp dir; set MigrationsDir + DSN + args status.
2. Assert exit 0; stdout has no apply/blocked/deferred data rows for fixtures.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	// intentionally empty — no writeMigration calls
	req.MigrationsDir = dir
	req.FixtureIDs = nil
	req.Args = []string{"status"}
	return nil
}
```
