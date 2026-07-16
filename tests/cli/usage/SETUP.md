# Scenario

**Feature**: DB subcommand usage errors when Config is incomplete

```
# missing cfg.DB (or MigrationsDir) for status/plan/apply/recovery → exit 2
cli.Run(cfg{DB:nil}, ["status"]) -> usage Error -> exit 2
```

## Preconditions

- Known subcommand that requires DB Config.
- No `--local`/`--remote` (those flags must not exist).
- Exit **2** for missing required Config fields (nil DB).

## Steps

1. Leaf sets incomplete Config + known command args.
2. Assert exit 2 and error messaging.

## Context

- P1: missing **DB** (not DSN). Harness leaves `req.DSN` empty so
  `buildConfig` does not Wrap.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.CloseStdin = false
	if req.Args == nil {
		req.Args = []string{}
	}
	t.Log("usage branch: incomplete Config (nil DB) → exit 2")
	return nil
}
```
