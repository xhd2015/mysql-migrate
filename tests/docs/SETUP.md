# Scenario

**Feature**: module root README documents purpose, install, CLI, env, and doctests

```
# operator / contributor reads module README
module root (go.mod, README.md)
  -> README.md documents purpose, install, CLI flags, subcommands, env, doctests

# this tree only reads the file (no CLI, no MySQL)
tests/docs Run -> os.ReadFile(README.md) -> phrase Assert per section leaf
```

## Preconditions

- Working directory for the product is the mysql-migrate repo root
  (`DOCTEST_ROOT/../..` from this tree).
- Module path in `go.mod` is `github.com/xhd2015/mysql-migrate`.
- Production surface under test is **`README.md` at the module root**
  (Classic RED until implementer adds it with the locked phrases).
- Out of scope: CLI behavior, MySQL, package APIs (covered by other trees).
- Do not import lifelog packages from this tree.

## Steps

1. Root setup validates the request pointer and that the module root path
   resolves (directory exists).
2. Leaves under `readme-section/` set `req.Label` and `req.RequiredPhrases`.
3. Root `Run` reads `README.md` from the module root and returns content or
   `Exists=false` when the file is missing.

## Context

- Locked README topics: purpose, install, CLI `--dsn`/`--dir`, seven
  subcommands, `MIGRATE_MYSQL_*` env vars, core DSN-free / `sqlexec`, how to
  run doctests.
- Phrase checks are case-sensitive substrings (implementer may use any
  surrounding prose as long as the tokens appear).
- Optional root DOCTEST index is out of scope for this tree (implementer).

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
	moduleRoot := filepath.Clean(filepath.Join(DOCTEST_ROOT, "..", ".."))
	st, err := os.Stat(moduleRoot)
	if err != nil {
		return fmt.Errorf("module root %s: %w", moduleRoot, err)
	}
	if !st.IsDir() {
		return fmt.Errorf("module root is not a directory: %s", moduleRoot)
	}
	return nil
}
```
