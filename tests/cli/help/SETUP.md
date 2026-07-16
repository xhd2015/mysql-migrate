# Scenario

**Feature**: help paths print Usage and exit 0

```
# help never talks to MySQL; no --local/--remote
cli.Run(cfg, [-h | <cmd> -h]) -> stdout Usage (ProgramName) -> exit 0
```

## Preconditions

- Invocation is a help request (root or subcommand `-h`).
- Expected exit code for every leaf under this branch: **0**.
- Help text goes primarily to **stdout** (stderr may be empty).
- Config may only need ProgramName; DSN/MigrationsDir unused for help.

## Steps

1. Leaf sets concrete help args (and optional ProgramName).
2. Run CLI and assert exit 0 + Usage tokens.

## Context

- Sibling of `unknown/`, `usage/` (exit 2), and DB-backed branches.
- Help paths stay offline (no MySQL).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Help branch: leaves set Args to -h or <cmd> -h.
	req.CloseStdin = false
	if req.ProgramName == "" {
		req.ProgramName = "mysql-migrate"
	}
	return nil
}
```
