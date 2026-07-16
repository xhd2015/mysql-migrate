# Scenario

**Feature**: apply stops on first SQL execution failure

```
# Action==apply starts; db.Exec fails → MarkFailed; later applyables not run
plan clear (no pre-block)
  -> apply first (bad SQL) fails
  -> stop exit 1; later pending untouched
```

## Preconditions

- Plan HasBlock is false at start.
- First migration body is invalid SQL that MySQL rejects.
- Later migration is a valid CREATE TABLE that must **not** run after failure.

## Steps

1. Leaves write bad-then-good fixtures.
2. Run apply.
3. Assert first failed in log, later not applied, exit 1.

## Context

- Distinct from refuse-block (no Exec at all when already blocked).

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
