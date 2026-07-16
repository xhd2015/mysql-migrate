# mysql-migrate — cmd binary `cmd/mysql-migrate` (P6)

Thin main for the operator tool. Parses **global** flags with **less-flags**,
builds `migrate.Config`, and delegates to **`cli.Run(cfg, remainArgs)`**.
Classic RED until implementer wires main beyond the empty stub.

Standalone doctest root under `tests/cmd/` (does not inherit `tests/cli/`).
CLI library behavior is sealed under `tests/cli/`; this tree only locks the
**binary entry surface**: global flags, env fallbacks, help, and hand-off.

Target:

```text
go run ./cmd/mysql-migrate [global flags] <command> [args]
# or built binary: mysql-migrate …
```

Implementer surface (suggested):

```go
// cmd/mysql-migrate
func main() { os.Exit(run(os.Args[1:])) }

// less-flags: --dsn, --dir (-h/--help), StopOnFirstArg
// empty remain / help → print root usage, exit 0
// else: Config{DSN, MigrationsDir, ProgramName} → cli.Run(cfg, remain)
```

- Global flags: `--dsn`, `--dir` (migrations directory).
- Env fallbacks (optional): `MIGRATE_MYSQL_DSN`, `MIGRATE_MYSQL_DIR`
  (flag wins when both set).
- Empty argv → **help, exit 0** (unlike bare `cli.Run` with empty args).
- Subcommand help: remaining args include e.g. `status -h` → `cli.Run` help.
- Missing DSN for DB subcommands (no flag, no env) → non-zero usage (exit **2**).

# DSN (Domain Specific Notion)

The **mysql-migrate binary** is a thin **operator front door**. An **operator**
runs the process with **argv** (and optional **environment**). The binary
**parses global flags** (`--dsn`, `--dir`) with **less-flags**, applies **env
fallbacks** when a flag is omitted, then **builds Config** (DSN, MigrationsDir,
ProgramName=`mysql-migrate`) and **calls `cli.Run`** with the remaining
subcommand args.

**Help** at the binary root (`-h` / `--help`, or **empty args**) prints
**Usage** listing subcommands and global flags, then exits **0** without
requiring DSN. **Subcommand help** (e.g. `status -h`) is handled after
global parse by **`cli.Run`**, also exit **0**.

**Usage errors** at the config boundary (e.g. `apply` with no DSN from flag or
env) surface as non-zero exit (**2**) with an **Error** about missing DSN —
without opening MySQL.

**Apply** (and other DB commands) open the DSN, use the migrations **dir**, and
follow the sealed **`cli.Run`** contracts (progress, log rows, exit codes).
This tree exercises one happy **apply** path through the binary flags to prove
wiring end-to-end when MySQL is available; it does not re-seal full CLI
coverage.

Tests **build** `./cmd/mysql-migrate` once per `doctest test` session (flock +
ready marker under `$TMPDIR`), **exec** the binary with controlled env (strip
ambient `MIGRATE_MYSQL_*` by default), and assert exit codes + stdout/stderr.

## Version

0.0.2

## Decision Tree

```
tests/cmd/                                   [Request{Args, Env, ClearMigrateEnv, …}]
│                                            Run: exec session-built mysql-migrate
├── help/                                    # offline, exit 0
│   ├── root/                                # -h → Usage lists commands + --dsn/--dir
│   ├── empty-args/                          # no args → help exit 0 (binary rule)
│   └── status/                              # status -h via binary → Usage status
├── usage/
│   └── apply-missing-dsn/                   # apply, no --dsn, no env → exit 2
└── apply/
    └── one-create-table/                    # --dsn + --dir apply; skip if MySQL down
```

**Significance order:** dispatch class (help | usage | apply) → help variant
(root / empty / subcommand) → config source (missing vs flags) → fixture apply.

## Test Index

| Leaf | Description |
|------|-------------|
| `help/root` | `-h` → exit 0; Usage lists all subcommands; mentions `--dsn` and `--dir` |
| `help/empty-args` | no args → exit 0; Usage (same root help contract) |
| `help/status` | `status -h` → exit 0; Usage mentions `status`; no `--local`/`--remote` |
| `usage/apply-missing-dsn` | `apply --dir <tmp>` without DSN flag/env → exit **2**; Error mentions dsn |
| `apply/one-create-table` | `--dsn` + `--dir` + `apply` one CREATE TABLE → exit **0**, log success; skip if MySQL unreachable |

## How to Run

```sh
cd /Users/xhd2015/Projects/xhd2015/mysql-migrate
# optional for apply leaf:
# export MIGRATE_MYSQL_DSN='user:pass@tcp(host:port)/db?multiStatements=true&...'
doctest vet ./tests/cmd
doctest test ./tests/cmd
```

Offline leaves: all `help/*` and `usage/apply-missing-dsn` (no MySQL).
`apply/one-create-table` **skips** when the harness DSN is not reachable.

Classic TDD: `cmd/mysql-migrate` is an empty `main` until implementer lands
less-flags + Config + `cli.Run`. Leaves must fail (empty output / wrong exit)
until the binary wires correctly.

### Binary surface (locked)

```text
mysql-migrate -h
mysql-migrate                          # empty → help exit 0
mysql-migrate status -h
mysql-migrate --dir <path> apply       # missing DSN → exit 2
mysql-migrate --dsn <dsn> --dir <path> apply
```

Env (optional, flag overrides):

```text
MIGRATE_MYSQL_DSN
MIGRATE_MYSQL_DIR
```

### Exit codes (binary / cli hand-off)

| Case | Exit |
|------|------|
| Root help / empty args / subcommand help | **0** |
| Apply success | **0** |
| Missing DSN (usage) | **2** |

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Request drives one invocation of the built mysql-migrate binary.
type Request struct {
	// Args are process argv after the program name (passed to the binary).
	Args []string

	// ClearMigrateEnv when true strips MIGRATE_MYSQL_DSN and MIGRATE_MYSQL_DIR
	// from the child environment before applying Env. Default true (root Setup).
	ClearMigrateEnv bool

	// Env is extra KEY=value pairs appended to the child environment.
	Env []string

	// Bin is the absolute path to the session-built mysql-migrate binary.
	// Root Setup fills this via buildBinaryOnce.
	Bin string

	// WorkDir is cmd.Dir; empty means module root.
	WorkDir string

	// FixtureIDs are migration_id values written by apply leaf Setup.
	FixtureIDs []string

	// TableNames are optional MySQL tables created by apply fixtures.
	TableNames []string

	// AssertDSN is the DSN used by Assert for log/table side effects
	// (same as the DSN passed to the binary for apply leaves).
	AssertDSN string

	// MigrationsDir is the temp migrations directory (apply fixtures).
	MigrationsDir string
}

// Response is the captured subprocess outcome.
type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
	// Duration is wall time of the subprocess (hang detection if needed).
	Duration time.Duration
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	if strings.TrimSpace(req.Bin) == "" {
		return nil, fmt.Errorf("Request.Bin empty — root Setup must build binary")
	}
	if req.Args == nil {
		req.Args = []string{}
	}

	workDir := req.WorkDir
	if workDir == "" {
		workDir = repoRoot(t)
	}

	cmd := exec.Command(req.Bin, req.Args...)
	cmd.Dir = workDir
	cmd.Env = buildChildEnv(req)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	runErr := cmd.Run()
	dur := time.Since(start)

	exitCode := 0
	if runErr != nil {
		if ee, ok := runErr.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			return nil, fmt.Errorf("exec %s: %w", req.Bin, runErr)
		}
	}

	return &Response{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: dur,
	}, nil
}

// buildChildEnv constructs the child process environment.
func buildChildEnv(req *Request) []string {
	base := os.Environ()
	if req.ClearMigrateEnv {
		filtered := make([]string, 0, len(base))
		for _, e := range base {
			// strip known migrate env keys (prefix match on KEY=)
			if strings.HasPrefix(e, "MIGRATE_MYSQL_DSN=") ||
				strings.HasPrefix(e, "MIGRATE_MYSQL_DIR=") {
				continue
			}
			filtered = append(filtered, e)
		}
		base = filtered
	}
	if len(req.Env) == 0 {
		return base
	}
	return append(append([]string{}, base...), req.Env...)
}

func repoRoot(t *testing.T) string {
	t.Helper()
	// DOCTEST_ROOT is tests/cmd → module root is ../..
	root, err := filepath.Abs(filepath.Join(DOCTEST_ROOT, "..", ".."))
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	return root
}

func sessionCacheDir() string {
	return filepath.Join(os.TempDir(), "mysql-migrate-cmd-doctest-"+DOCTEST_SESSION_ID)
}

func withFileLock(t *testing.T, lockPath string, fn func() error) error {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return fn()
}

// buildBinaryOnce builds ./cmd/mysql-migrate once per doctest session (flock).
func buildBinaryOnce(t *testing.T) string {
	t.Helper()
	cache := sessionCacheDir()
	lock := filepath.Join(cache, "build.lock")
	ready := filepath.Join(cache, "binaries.ready")
	bin := filepath.Join(cache, "mysql-migrate")
	root := repoRoot(t)

	err := withFileLock(t, lock, func() error {
		if st, e := os.Stat(ready); e == nil && st.Mode().IsRegular() {
			if st2, e2 := os.Stat(bin); e2 == nil && st2.Mode().IsRegular() {
				return nil
			}
		}
		if err := os.MkdirAll(cache, 0o755); err != nil {
			return err
		}
		cmd := exec.Command("go", "build", "-o", bin, "./cmd/mysql-migrate")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("go build ./cmd/mysql-migrate: %w\n%s", err, strings.TrimSpace(string(out)))
		}
		return os.WriteFile(ready, []byte("ok"), 0o644)
	})
	if err != nil {
		t.Fatalf("buildBinaryOnce: %v", err)
	}
	return bin
}
```
