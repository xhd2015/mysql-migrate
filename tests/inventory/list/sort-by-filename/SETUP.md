# Scenario

**Feature**: ListDir returns migrations sorted by full filename ascending

```
# write files out of chronological order; list must sort by name
dir has:
  2026-07-17-01-b.sql
  2026-07-16-02-a.sql
  2026-07-16-01-z.sql
ListDir -> [01-z, 02-a, 17-01-b] by FileName ASC
```

## Preconditions

- Three valid top-level migration files written with unordered creation.

## Steps

1. Write three valid `.sql` fixtures with non-sorted names.
2. ListDir and expect FileName order:
   - `2026-07-16-01-z.sql`
   - `2026-07-16-02-a.sql`
   - `2026-07-17-01-b.sql`

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "list"
	req.FixtureFiles = map[string]string{
		// Intentionally out of order map insertion / names
		"2026-07-17-01-b.sql": "-- b\n",
		"2026-07-16-02-a.sql": "-- a\n",
		"2026-07-16-01-z.sql": "-- z\n",
	}
	return nil
}
```
