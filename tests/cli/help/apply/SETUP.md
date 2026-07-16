# Scenario

**Feature**: `apply -h` prints apply Usage and exits 0

```
# apply help surfaces optional --to (no --local/--remote)
cli.Run(cfg, ["apply", "-h"]) -> apply Usage -> exit 0
```

## Preconditions

- Args: `["apply", "-h"]`.
- Help must not require DSN/flags to print.

## Steps

1. Set args to apply help.
2. Assert exit 0 and Usage mentions apply and `--to`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"apply", "-h"}
	req.ProgramName = "mysql-migrate"
	return nil
}
```
