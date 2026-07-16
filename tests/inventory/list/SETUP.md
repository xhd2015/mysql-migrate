# Scenario

**Feature**: ListDir inventories top-level matching migration SQL files

```
# list top-level only, sort by filename, fill ContentSHA256
ListDir(dir) -> []MigrationFile | error
# non-matching junk ignored; grammar lookalikes error
```

## Preconditions

- `Dir` is a temporary directory prepared with fixture files.
- Listing is **not recursive** — nested `.sql` under subdirs is ignored.
- Sort order is lexicographic on full filename (= date then NN).

## Steps

1. Set `req.Op = "list"`.
2. Create temp `req.Dir` and populate `FixtureFiles` / `FixtureNested` in leaves.
3. Run `ListDir`.

## Context

- Prefer: names matching loose `YYYY-MM-DD-*.sql` that fail full grammar → **error**.
- Plain unrelated files (`README`, `.gitkeep`, `notes.txt`) → **ignored**, not error.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "list"
	req.Dir = newTempDir(t)
	return nil
}
```
