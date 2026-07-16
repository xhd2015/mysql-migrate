# Scenario

**Feature**: root `-h` prints full command Usage and global flags, exit 0

```
# operator asks for top-level help
mysql-migrate -h
  -> Usage listing subcommands + --dsn/--dir -> exit 0
```

## Preconditions

- Args: `["-h"]`.
- Root help must list every locked subcommand name.
- Root help must document global `--dsn` and `--dir`.

## Steps

1. Set `req.Args` to `[]string{"-h"}`.
2. Run binary.
3. Assert exit 0, command names, and global flags on stdout.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-h"}
	return nil
}
```
