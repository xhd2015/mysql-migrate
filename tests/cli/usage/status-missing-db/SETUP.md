# Scenario

**Feature**: `status` with nil `cfg.DB` is a usage error (exit 2)

```
# DB comes from Config, not flags; harness does not Wrap when DSN empty
cli.Run(cfg{DB:nil, MigrationsDir:tmp}, ["status"])
  -> exit 2; Error about DB / missing config
```

## Preconditions

- Args: `["status"]` (no target flags).
- `req.DSN` empty string so harness leaves `cfg.DB` nil.
- `req.MigrationsDir` may be a temp dir (failure must be DB, not dir).
- Offline — no MySQL required (fails before any open).

## Steps

1. Create temp migrations dir; leave harness DSN empty.
2. Run status.
3. Expect exit 2; combined output mentions db/config/missing.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	req.MigrationsDir = dir
	req.DSN = "" // harness: do not open/wrap → cfg.DB nil
	req.Args = []string{"status"}
	t.Log("status-missing-db: intentionally nil cfg.DB")
	return nil
}
```
