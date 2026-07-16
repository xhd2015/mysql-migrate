# Scenario

**Feature**: valid migration basenames parse successfully

```
# well-formed name yields metadata without error
ParseFileName(valid) -> MigrationFile (err == nil)
```

## Preconditions

- Names under this branch conform to the full grammar.

## Steps

1. Keep `req.Op = "parse"`.
2. Leaf sets a concrete valid `FileName`.

## Context

- Cover simple slug, middle `[EXACTLY-ONCE]`, and stem without `.sql`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "parse"
	return nil
}
```
