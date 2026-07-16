# Scenario

**Feature**: DB subcommand usage errors when Config is incomplete

```
# missing cfg.DSN (or MigrationsDir) for status/plan/apply/recovery → exit 2
cli.Run(cfg{DSN:""}, ["status"]) -> usage Error -> exit 2
```

## Preconditions

- Known subcommand that requires DB Config.
- No `--local`/`--remote` (those flags must not exist).
- Exit **2** for missing required Config fields.

## Steps

1. Leaf sets incomplete Config + known command args.
2. Assert exit 2 and error messaging.

## Context

- Replaces lifelog `missing-target` (`--local`/`--remote`) with **missing DSN**.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.CloseStdin = false
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
