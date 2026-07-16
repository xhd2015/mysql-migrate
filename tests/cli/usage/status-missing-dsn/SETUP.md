# Scenario

**Feature**: `status` with empty `cfg.DSN` is a usage error (exit 2)

```
# DSN comes from Config, not flags
cli.Run(cfg{DSN:"", MigrationsDir:tmp}, ["status"])
  -> exit 2; Error about DSN / missing config
```

## Preconditions

- Args: `["status"]` (no target flags).
- `req.DSN` empty string.
- `req.MigrationsDir` may be a temp dir (failure must be DSN, not dir).
- Offline — no MySQL required (fails before open).

## Steps

1. Create temp migrations dir; leave DSN empty.
2. Run status.
3. Expect exit 2; combined output mentions dsn/config/missing.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	req.MigrationsDir = dir
	req.DSN = "" // intentionally empty
	req.Args = []string{"status"}
	return nil
}
```
