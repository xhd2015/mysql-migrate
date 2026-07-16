# Scenario

**Feature**: `status -h` via binary prints status Usage and exits 0

```
# global parse leaves remain [status, -h] for cli.Run
mysql-migrate status -h
  -> Usage for status -> exit 0 (offline)
```

## Preconditions

- Args: `["status", "-h"]`.
- No global `--dsn` / `--dir` required for help.
- Status help must not document removed `--local` / `--remote` flags.

## Steps

1. Set args to `status -h`.
2. Run binary.
3. Assert exit 0 and status Usage tokens.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"status", "-h"}
	return nil
}
```
