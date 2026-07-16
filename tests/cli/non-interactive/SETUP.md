# Scenario

**Feature**: CLI never blocks waiting for interactive stdin

```
# closed or empty stdin must not hang
closed stdin -> cli.Run(cfg, ...) finishes promptly
```

## Preconditions

- Stdin provides immediate EOF (closed pipe) or is `/dev/null`.
- Operator tools must not prompt for passwords/SSH.

## Steps

1. Leaf sets `CloseStdin` and a safe args set (help).
2. Assert Run returns quickly with expected exit.

## Context

- Complements usage paths; focuses on wall-clock duration under closed stdin.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.CloseStdin = true
	return nil
}
```
