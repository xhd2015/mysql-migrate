# Scenario

**Feature**: apply refuses when plan already HasBlock (no later Exec)

```
# blocked item before applyables → refuse; do not apply deferred later migrations
plan.Build HasBlock=true
  -> stderr Error … blocked …
  -> exit 1; no MarkRunning/Exec for deferred items
```

## Preconditions

- First migration is in a **blocked** state (e.g. EXACTLY-ONCE **failed**).
- Later migration would be **deferred** by plan stop-chain.
- Apply must not Exec later SQL.

## Steps

1. Leaves seed block condition + write later pending fixture.
2. Run apply.
3. Assert exit 1, Error+blocked on stderr, later not applied.

## Context

- Distinct from exec-fail (which starts apply then fails SQL).

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
