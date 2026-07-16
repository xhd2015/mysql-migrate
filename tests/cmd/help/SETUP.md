# Scenario

**Feature**: binary help paths print Usage and exit 0 (offline)

```
# root help is binary-owned; subcommand help after global parse
mysql-migrate -h | (no args) | status -h
  -> Usage on stdout -> exit 0 (no MySQL)
```

## Preconditions

- Invocation is a help request (root `-h`, empty args, or `status -h`).
- Expected exit code for every leaf under this branch: **0**.
- Help text primarily on **stdout**.
- No DSN or migrations dir required.
- Child env has migrate vars stripped (root default).

## Steps

1. Leaf sets concrete Args for the help form under test.
2. Exec binary; Assert exit 0 + Usage tokens.

## Context

- Sibling of `usage/` (exit 2) and `apply/` (DB).
- Empty-args help is **cmd-specific** (library `cli.Run` treats empty as usage).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Help is offline and must not inherit ambient DSN/dir env.
	req.ClearMigrateEnv = true
	req.Env = nil
	return nil
}
```
