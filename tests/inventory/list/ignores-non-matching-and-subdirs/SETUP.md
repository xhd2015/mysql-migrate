# Scenario

**Feature**: ListDir ignores non-matching top-level entries and all nested files

```
# only valid top-level grammar matches appear
dir:
  2026-07-16-01-keep-me.sql          # included
  README                              # ignored
  notes.txt                           # ignored
  nested/2026-07-16-99-nested.sql     # ignored (not top-level)
ListDir -> [2026-07-16-01-keep-me.sql] only
```

## Preconditions

- One valid migration at top level.
- Junk files and a nested valid-looking `.sql` under a subdirectory.

## Steps

1. Write valid top-level migration, junk files, and nested sql.
2. ListDir must return only the top-level valid file.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "list"
	req.FixtureFiles = map[string]string{
		"2026-07-16-01-keep-me.sql": "-- keep\n",
		"README":                    "not a migration\n",
		"notes.txt":                 "ignore me\n",
		".gitkeep":                  "",
	}
	req.FixtureNested = map[string]string{
		"nested/2026-07-16-99-nested.sql": "-- nested should not list\n",
		"migrate_legacy/tool.sql":         "-- legacy tool path\n",
	}
	return nil
}
```
