# Scenario

**Feature**: parse a simple migration name without EXACTLY-ONCE

```
# standard append-only migration filename
ParseFileName("2026-07-16-01-create-t-channel.sql")
  -> ID=2026-07-16-01-create-t-channel
  -> Date=2026-07-16, Seq=1, ExactlyOnce=false, Slug=create-t-channel
```

## Preconditions

- Basename: `2026-07-16-01-create-t-channel.sql`

## Steps

1. Set `FileName` to the simple valid example.
2. Run `ParseFileName`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "parse"
	req.FileName = "2026-07-16-01-create-t-channel.sql"
	return nil
}
```
