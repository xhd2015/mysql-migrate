# Scenario

**Feature**: empty argv shows root help and exits 0 (binary rule)

```
# no subcommand, no flags — binary treats as help, not usage error
mysql-migrate
  -> same root Usage as -h -> exit 0
```

## Preconditions

- Args: empty slice `[]string{}` (no tokens after program name).
- Must **not** exit 2 (that would be bare `cli.Run` empty-args behavior).
- Same Usage content expectations as `help/root`.

## Steps

1. Set `req.Args` to empty.
2. Run binary.
3. Assert exit 0 + root help tokens.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{}
	return nil
}
```
