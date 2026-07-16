# Scenario

**Feature**: `status` builds full plan table against DB + migrations dir (Config)

```
# status shows every migration item (apply/skip/blocked/deferred)
cfg{DSN, MigrationsDir=<fixture>}
cli.Run(cfg, ["status"])
  -> EnsureTable + ListDir + List + plan.Build
  -> stdout all items; exit 0 if !HasBlock else 1
```

## Preconditions

- Args always include `status` (no `--local`/`--remote`).
- `req.MigrationsDir` + `req.DSN` set via `fillConfigForDB` / leaf fixtures.
- Local MySQL reachable (`ensureMySQL`); skip if down.
- Leaves seed log rows only for their fixture migration_ids.

## Steps

1. Group Setup ensures MySQL + fills DSN; leaves write fixtures / seed logs / set Args.
2. Run status with Config.
3. Assert exit code + action tokens for FixtureIDs.

## Context

- Sibling of `plan/` (plan filters to non-skip).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	fillConfigForDB(t, req)
	req.CloseStdin = false
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
