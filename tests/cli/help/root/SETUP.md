# Scenario

**Feature**: root `-h` prints full command Usage and exits 0

```
# operator asks for top-level help
cli.Run(cfg, ["-h"]) -> Usage listing all subcommands -> exit 0
```

## Preconditions

- Args: `["-h"]`.
- Root help must list every locked subcommand name.
- Usage should mention ProgramName from Config.

## Steps

1. Set `req.Args` to `[]string{"-h"}` and a known ProgramName.
2. Run CLI.
3. Assert exit 0 and command names present on stdout.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-h"}
	req.ProgramName = "mysql-migrate"
	return nil
}
```
