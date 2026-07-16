# Scenario

**Feature**: ListDir errors on top-level names that look like migrations but fail grammar

```
# loose YYYY-MM-DD-*.sql that is not full grammar -> ListDir error
dir:
  2026-07-16-01-ok.sql           # valid
  2026-07-16-1-bad-seq.sql       # looks like migration, NN not padded
ListDir -> error (do not silently skip lookalikes)
```

## Preconditions

- Prefer policy from requirement: only consider names matching loose
  `YYYY-MM-DD-*.sql`; if grammar fails, **error** (not silent ignore).
- Plain junk without that loose prefix remains ignored (covered elsewhere).

## Steps

1. Write one valid file and one lookalike with unpadded `NN`.
2. ListDir must return a non-nil error.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "list"
	req.FixtureFiles = map[string]string{
		"2026-07-16-01-ok.sql":     "-- ok\n",
		"2026-07-16-1-bad-seq.sql": "-- lookalike invalid NN\n",
	}
	return nil
}
```
