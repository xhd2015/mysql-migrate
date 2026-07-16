# Scenario

**Feature**: `status -h` prints status Usage and exits 0

```
# subcommand help (no target flags; DSN is from Config)
cli.Run(cfg, ["status", "-h"]) -> status Usage -> exit 0
```

## Preconditions

- Args: `["status", "-h"]`.
- Must not require DSN when asking for help.

## Steps

1. Set args to status help.
2. Assert exit 0 and Usage mentions status.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"status", "-h"}
	req.ProgramName = "mysql-migrate"
	return nil
}
```
