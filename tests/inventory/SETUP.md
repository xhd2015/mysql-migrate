# Scenario

**Feature**: pure migration file inventory — parse names, list dirs, hash bytes

```
# author places top-level migration SQL files under migrations root
author -> migrations dir (YYYY-MM-DD-NN[-[EXACTLY-ONCE]]-<slug>.sql)

# inventory parses basename metadata (no DB)
ParseFileName(basename) -> MigrationFile{ID, Date, Seq, ExactlyOnce, Slug}

# list scans top-level only, sorts by filename, fills ContentSHA256
ListDir(dir) -> []MigrationFile (lexicographic FileName ASC)

# hash is raw-byte SHA-256 lowercase hex
HashFile(path) -> 64-char hex digest
```

## Preconditions

- Module path is `github.com/xhd2015/mysql-migrate` (repo `go.mod`).
- Target package path (Classic RED until implemented):
  `migrate/inventory`
  import `github.com/xhd2015/mysql-migrate/migrate/inventory`
- Filename grammar:
  `YYYY-MM-DD-NN[-[EXACTLY-ONCE]]-<slug>.sql`
  - `NN` zero-padded `01`–`99`
  - optional exact middle token `[EXACTLY-ONCE]`
  - `slug` non-empty kebab-case
  - `migration_id` = basename without `.sql`
- Hash algorithm: SHA-256 of raw file bytes, lowercase hex (64 chars).
- No MySQL, no CLI, no apply log in this inventory tree (P2).
- Out of scope: plan, logrepo, CLI, lifelog.

## Steps

1. Leaf Setup sets `req.Op` and either a basename, an invalid-name table, or
   fixture file maps under a temp `req.Dir`.
2. Root helper materialises fixture maps into the temp directory when present.
3. `Run` dispatches to `inventory.ParseFileName`, `inventory.ListDir`, or
   `inventory.HashFile`.
4. Leaf `Assert` checks metadata, sort order, ignore rules, or digests.

## Context

- Package is stub-only until implementer ports logic — tests stay Classic RED.
- Fixtures use per-leaf temp dirs — no shared mutable migrations directory.
- Parallel-safe: no global filesystem under a real migrations path.
- Helpers `materializeFixtures` and `newTempDir` are shared by list/hash leaves.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	if req == nil {
		return fmt.Errorf("nil request")
	}
	// Ensure maps exist so leaf Setups can assign without nil panics.
	if req.FixtureFiles == nil {
		req.FixtureFiles = map[string]string{}
	}
	if req.FixtureNested == nil {
		req.FixtureNested = map[string]string{}
	}
	return nil
}

// materializeFixtures writes top-level and nested relative paths under dir.
func materializeFixtures(t *testing.T, dir string, top map[string]string, nested map[string]string) error {
	t.Helper()
	if dir == "" {
		return fmt.Errorf("materializeFixtures: empty dir")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for name, content := range top {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	for rel, content := range nested {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// newTempDir returns an empty temp directory cleaned up with t.Cleanup.
func newTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "mysql-migrate-inventory-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}
```
