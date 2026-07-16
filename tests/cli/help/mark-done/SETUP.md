# Scenario

**Feature**: `mark-done -h` prints Usage mentioning `--note` and exits 0

```
# recovery subcommand help documents required --note
cli.Run(cfg, ["mark-done", "-h"]) -> Usage with --note -> exit 0
```

## Preconditions

- Args: `["mark-done", "-h"]`.
- Recovery commands require `--note` when run; help must surface that flag.

## Steps

1. Set args to mark-done help.
2. Assert exit 0 and stdout contains `--note`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"mark-done", "-h"}
	req.ProgramName = "mysql-migrate"
	return nil
}
```
