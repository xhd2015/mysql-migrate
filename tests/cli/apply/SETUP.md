# Scenario

**Feature**: `apply` mutates DB per plan (Config DSN + MigrationsDir)

```
# apply builds the same plan as status, then executes Action==apply items
cfg{DSN, MigrationsDir=<fixture>}
cli.Run(cfg, ["apply", ...optional --to])
  -> sql.Open -> EnsureTable -> ListDir -> List -> plan.Build
  -> refuse if blocked, else MarkRunning -> Exec -> MarkSuccess|MarkFailed
  -> progress + summary on stdout; exit 0|1
```

## Preconditions

- Args always include `apply` (no `--local`/`--remote`).
- `req.MigrationsDir` + DSN set; DSN should allow multiStatements.
- Local MySQL reachable; skip if down.
- Fixture migration_ids and optional table names are unique per session/leaf.

## Steps

1. Group Setup ensures MySQL; leaves write SQL fixtures / seed logs / set Args.
2. Run apply (optional `--to`).
3. Assert exit code, progress/summary tokens, log status, and table side effects.

## Context

- Split children by plan outcome: **clear** | **exec-fail** | **refuse-block**.

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
