# Scenario

**Feature**: Config is DB-only; nil DB is a CLI usage error

```
# Config shape
migrate.Config{DB sqlexec.DB, MigrationsDir, ProgramName, AppliedBy}
  // no DSN field

# nil DB usage
cli.Run(cfg{DB:nil, MigrationsDir:tmp}, ["status"]) -> exit 2
```

## Preconditions

- Offline: no MySQL required for either leaf.
- Depends on `migrate` Config field change and `cli` requiring `cfg.DB`.

## Steps

1. Leaves set Op to `config_fields` or `cli_nil_db`.
2. Assert field shape or exit code 2.

## Context

- Sibling surface to sqlexec methods; locks the engine injection contract.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Config branch is offline; leaves pick Op.
	t.Log("config branch: Config.DB-only + nil DB usage (offline)")
	return nil
}
```
