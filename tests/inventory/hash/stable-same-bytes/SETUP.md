# Scenario

**Feature**: HashFile is stable for the same file bytes

```
# two HashFile calls on the same path return identical digests
write body "hello-migrate\n"
HashFile(path) == HashFile(path)
```

## Preconditions

- Single file with fixed content.

## Steps

1. Write `stable.sql` with fixed body under temp dir.
2. Set `Path` only (no PathB) so Run hashes the path twice into Hash and HashB.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "hash"
	const body = "hello-migrate\n"
	req.Path = writeHashFile(t, req.Dir, "stable.sql", body)
	req.PathB = "" // Run re-hashes Path for HashB
	return nil
}
```
