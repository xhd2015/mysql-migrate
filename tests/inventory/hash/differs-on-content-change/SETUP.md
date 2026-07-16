# Scenario

**Feature**: different file contents produce different SHA-256 digests

```
# two files, different bodies
HashFile(a) != HashFile(b)
```

## Preconditions

- Two files under the same temp dir with distinct contents.

## Steps

1. Write `a.sql` with body A and `b.sql` with body B.
2. Set `Path` and `PathB`.
3. Run HashFile on both.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "hash"
	req.Path = writeHashFile(t, req.Dir, "a.sql", "content-A\n")
	req.PathB = writeHashFile(t, req.Dir, "b.sql", "content-B\n")
	return nil
}
```
