# Scenario

**Feature**: parse EXACTLY-ONCE marker in the middle of the filename

```
# optional middle token after NN-
ParseFileName("2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp.sql")
  -> ExactlyOnce=true, Slug=drop-legacy-tmp
  -> ID=2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp
```

## Preconditions

- Basename: `2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp.sql`
- Marker is the exact token `[EXACTLY-ONCE]` after `NN-`, before the slug.

## Steps

1. Set `FileName` to the EXACTLY-ONCE example.
2. Run `ParseFileName`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "parse"
	req.FileName = "2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp.sql"
	return nil
}
```
