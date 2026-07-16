# Scenario

**Feature**: unknown subcommand is a usage error (exit 2)

```
# bad first token is not silent
cli.Run(cfg, ["not-a-real-command"]) -> stderr Error -> exit 2
```

## Preconditions

- First arg is not a known subcommand and not a help flag.
- Exit code must be **2** (usage).

## Steps

1. Leaf sets a bogus subcommand name.
2. Assert exit 2 and stderr signals an error.

## Context

- Distinct from missing-DSN usage (known command, bad Config) and DB paths.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.CloseStdin = false
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
