# Scenario

**Feature**: root help with closed stdin completes quickly (no hang)

```
# stdin closed before Run; -h must not block
CloseStdin + cli.Run(cfg, ["-h"]) -> exit 0 within 2s
```

## Preconditions

- Args: `["-h"]` (stable, exit 0 path once implemented).
- `CloseStdin=true` — write end of stdin pipe closed before `cli.Run`.
- Hang budget: **2 seconds** wall time inside `cli.Run`.

## Steps

1. Set args to `-h` and force closed stdin.
2. Run CLI.
3. Assert exit 0 and `Duration < 2s`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-h"}
	req.CloseStdin = true
	req.ProgramName = "mysql-migrate"
	return nil
}
```
