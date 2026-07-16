# Scenario

**Feature**: completely unknown subcommand name exits 2 with Error on stderr

```
cli.Run(cfg, ["not-a-real-command"]) -> exit 2, stderr contains Error
```

## Preconditions

- Args: `["not-a-real-command"]` — not in the locked command set.
- Not a help flag (`-h` / `--help`).

## Steps

1. Set bogus args.
2. Run CLI.
3. Assert exit 2 and error-ish stderr.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"not-a-real-command"}
	return nil
}
```
