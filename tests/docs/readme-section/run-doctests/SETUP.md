# Scenario

**Feature**: README documents how to run the doctest suite

```
# run doctests section
README.md -> "doctest test" (and preferably tests/ path or go test)
```

## Preconditions

- README.md at module root (Classic RED if missing).
- Contributors need a copy-pasteable way to run the suite.

## Steps

1. Set `req.Label` to `run-doctests`.
2. Require phrases: `doctest test` and `doctest vet` so both validation and
   execution are documented.

## Context

- P7 exit criterion: document how to run doctests.
- `go test ./...` may also appear; not required by this leaf (doctest is the
  primary suite entry for this repo's design trees).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Label = "run-doctests"
	req.RequiredPhrases = []string{
		"doctest test",
		"doctest vet",
	}
	return nil
}
```
