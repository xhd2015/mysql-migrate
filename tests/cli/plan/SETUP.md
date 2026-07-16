# Scenario

**Feature**: `plan` shows non-skip actions (what would run / is gated)

```
# plan filters to apply | blocked | deferred (skip omitted)
cfg{DSN, MigrationsDir=<fixture>}
cli.Run(cfg, ["plan"])
  -> same pipeline as status
  -> stdout non-skip items; exit 0 if !HasBlock else 1
```

## Preconditions

- Args include `plan` (no target flags).
- Isolated temp migrations dir on Config.
- Local MySQL via `fillConfigForDB` / ensureMySQL.
- Output contract: **non-skip only** (apply, blocked, deferred).

## Steps

1. Group Setup ensures MySQL; leaves write fixtures / seed / set Args.
2. Run plan.
3. Assert exit + action tokens.

## Context

- Sibling of `status/` (status prints all items including skip).

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
