# mysql-migrate — CLI library `cli.Run(cfg, args)` (P5–P8)

Non-main reusable CLI tests for **status**, **plan**, **apply**, and human
**recovery** against an injected `migrate.Config` (DSN + MigrationsDir +
ProgramName + AppliedBy). Classic RED until implementer lands `cli.Run`.

Standalone doctest root under `tests/cli/` so inventory / plan / logrepo /
scaffold trees stay independent.

Target package (implementer provides):

```text
github.com/xhd2015/mysql-migrate/cli
```

Public API:

```go
// package cli
// Config is migrate.Config from github.com/xhd2015/mysql-migrate/migrate
func Run(cfg migrate.Config, args []string) int
```

- **Never** `os.Exit` — return exit codes only.
- **No** `--local` / `--remote` flags (DSN comes from `cfg.DSN`).
- DB subcommands require non-empty `cfg.DSN` and `cfg.MigrationsDir`.
- `cfg.ProgramName` appears in help `Usage` text (default-friendly when empty).
- Subcommands: `status`, `plan`, `apply [--to]`, `mark-done`, `mark-failed`,
  `note`, `allow-retry`.
- Exit **0** / **1** / **2** as lifelog design (help/success | biz | usage).

Pipeline (implementer):

```text
sql.Open(cfg.DSN) → logrepo.EnsureTable → inventory.ListDir(cfg.MigrationsDir)
  → logrepo.List → plan.Build
  → status/plan: print table (+ warnings)
  → apply: refuse if blocked; else for each Action==apply:
       MarkRunning → db.Exec(file SQL) → MarkSuccess | MarkFailed+stop
  → recovery: parse migration_id + required --note
       mark-done   → logrepo.MarkDone
       mark-failed → logrepo.MarkFailedManual
       note        → logrepo.SetNote
       allow-retry → logrepo.AllowRetry (EO only)
```

Thin main (not exercised by this tree):

```text
cmd/mysql-migrate/main.go → os.Exit(cli.Run(cfg, os.Args[1:]))
```

# DSN (Domain Specific Notion)

The **migrate CLI library** is a **non-interactive** operator tool. An
**operator** (or thin main) invokes **`cli.Run(cfg, args)`** with a
**Config** (DSN, migrations directory, program name, applied-by identity) and
**args** without the program name. The CLI **dispatches** on the first token:
**help** (`-h` / `--help` / empty+help), a **known subcommand** (`status`,
`plan`, `apply`, `mark-done`, `mark-failed`, `note`, `allow-retry`), or
**unknown** (usage error). Every command supports **help** that prints
**Usage** (using **ProgramName**) and exits **0**. Subcommands that talk to a
DB require non-empty **cfg.DSN** and **cfg.MigrationsDir**. Missing either is a
**usage** error (exit **2**). Unknown subcommands are usage errors (exit **2**)
with an **Error** line on stderr. There are **no** `--local` / `--remote`
flags — target selection is entirely via **Config**.

**Status** and **plan** are **read-mostly** operators: open DB with **cfg.DSN**,
**ensure** the log table, **list** migration files from **cfg.MigrationsDir**,
**list** log rows, **build** a plan via **plan.Build**, then **print** a table
to **stdout**. **Status** prints **all** plan items (id, action/status,
duration, note, exactly_once). **Plan** prints **non-skip** items only
(`apply`, `blocked`, `deferred`). Exit **0** when `HasBlock` is false; exit
**1** when any item is **blocked**. Hash mismatches emit a **`warning:`** line
on **stderr** (still exit 1 if blocked).

**Apply** is a **mutating** operator with the same plan pipeline, then:

1. If the plan has a **blocked** item — **refuse**: print blocked instructions
   (**Error** + **blocked** on stderr), exit **1**, do **not** run later
   deferred migrations.
2. Otherwise walk items with `Action==apply` in order:
   - optional **`--to <migration_id>`** stops after applying that id (inclusive)
   - **MarkRunning**(id, exactlyOnce, file hash, cfg.AppliedBy)
   - read file SQL; **db.Exec** (DSN should enable `multiStatements=true`)
   - success → **MarkSuccess** + duration; progress `apply <id> ... ok (Nms)`
   - error → **MarkFailed** + error; progress `... failed`; **STOP** exit **1**
3. Print a final summary: **N applied**, **N failed**, **N pending**.
4. Exit **0** when no failures.

**Recovery** is a **non-interactive human gate**. All four commands require
**exactly**: `<migration_id>` and **`--note` with non-empty value**. Missing
either → usage exit **2** (no prompt). Bodies call logrepo:

| Command | logrepo | Result |
|---------|---------|--------|
| `mark-done` | `MarkDone` | force **success** + note (no SQL re-run) |
| `mark-failed` | `MarkFailedManual` | force **failed** + note |
| `note` | `SetNote` | update note only; status unchanged |
| `allow-retry` | `AllowRetry` | EO failed → **pending** + note; non-EO → biz error exit **1** |

The tool must **never hang on stdin**. `Run` never calls `os.Exit`.

## Version

0.0.2

## Decision Tree

```
tests/cli/                                   [Request{Args, DSN, MigrationsDir, …}]
│                                            Run: cli.Run(cfg, args) [+ follow-up]
├── help/                                    # exit 0, Usage on stdout (offline)
│   ├── root/                                # -h lists all subcommands + ProgramName
│   ├── status/                              # status -h (no --local/--remote)
│   ├── apply/                               # apply -h surfaces --to
│   └── mark-done/                           # mark-done -h mentions --note
├── unknown/
│   └── subcommand/                          # bogus name → exit 2 + Error
├── usage/
│   └── status-missing-dsn/                  # status with empty cfg.DSN → 2
├── status/                                  # real status + fixture dir + DSN
│   ├── empty-migrations/                    # empty dir → exit 0
│   ├── all-pending/                         # fixtures, no logs → apply, exit 0
│   ├── with-success-log/                    # success log → skip (+ later apply)
│   └── hash-mismatch-warning/               # success hash ≠ file → blocked, exit 1
├── plan/                                    # real plan + fixture dir + DSN
│   ├── all-pending/                         # non-skip apply rows, exit 0
│   └── eo-failed-blocked/                   # EO failed → blocked + deferred, exit 1
├── apply/                                   # real apply + fixture dir + DSN
│   ├── clear/                               # plan HasBlock=false
│   │   ├── two-create-tables/               # two CREATE TABLE → success
│   │   ├── second-apply-all-skip/           # prior success → all skip, exit 0
│   │   └── to-mid/                          # --to mid id → only up to that id
│   ├── exec-fail/
│   │   └── bad-sql-stops-later/             # bad SQL → failed; later not applied
│   └── refuse-block/
│       └── eo-failed/                       # EO failed → exit 1; later not applied
├── recovery/                                # mark-* / note / allow-retry
│   ├── mark-done/
│   │   ├── then-status-skip/                # mark-done → success; status skip
│   │   └── empty-note/                      # --note "" → exit 2
│   ├── mark-failed/
│   │   └── with-note/                       # mark-failed → failed + note
│   ├── note/
│   │   └── update-only/                     # note → note updated; status same
│   ├── allow-retry/
│   │   ├── eo-then-apply/                   # EO failed → allow-retry → apply ok
│   │   └── non-eo/                          # non-EO → exit 1 Error
│   └── usage/
│       ├── missing-note/                    # no --note → exit 2
│       └── missing-id/                      # no migration_id → exit 2
└── non-interactive/
    └── closed-stdin/                        # closed stdin + -h finishes quickly
```

**Significance order:** dispatch class (help | unknown | usage | status | plan |
apply | recovery | non-interactive) → recovery command / apply outcome class →
happy vs error variant.

## Test Index

| Leaf | Description |
|------|-------------|
| `help/root` | `-h` → exit 0; Usage lists all subcommands; mentions ProgramName |
| `help/status` | `status -h` → exit 0; Usage for status (no `--local`/`--remote`) |
| `help/apply` | `apply -h` → exit 0; mentions `--to` |
| `help/mark-done` | `mark-done -h` → exit 0; mentions `--note` |
| `unknown/subcommand` | Unknown name → exit 2; stderr contains Error |
| `usage/status-missing-dsn` | `status` with empty `cfg.DSN` → exit 2 |
| `status/empty-migrations` | Empty migrations dir → exit 0 |
| `status/all-pending` | Two fixtures, no logs → both **apply**, exit 0 |
| `status/with-success-log` | First success → **skip**; second **apply**; exit 0 |
| `status/hash-mismatch-warning` | Hash mismatch → **blocked**, exit **1**, `warning:` |
| `plan/all-pending` | Plan shows **apply** rows, exit 0 |
| `plan/eo-failed-blocked` | EO failed + later → **blocked** + **deferred**, exit **1** |
| `apply/clear/two-create-tables` | Two CREATE TABLE → success logs + tables, exit **0** |
| `apply/clear/second-apply-all-skip` | Prior success → no re-fail, exit **0** |
| `apply/clear/to-mid` | `--to` mid → first two applied; third pending |
| `apply/exec-fail/bad-sql-stops-later` | Bad SQL → first **failed**, later not applied, exit **1** |
| `apply/refuse-block/eo-failed` | EO failed seed → exit **1**, later not applied |
| `recovery/mark-done/then-status-skip` | mark-done → success+note; follow-up status **skip** |
| `recovery/mark-done/empty-note` | `--note ""` → exit **2** |
| `recovery/mark-failed/with-note` | mark-failed → log **failed** + note, exit **0** |
| `recovery/note/update-only` | note → note updated, status still **success** |
| `recovery/allow-retry/eo-then-apply` | allow-retry → follow-up apply success once |
| `recovery/allow-retry/non-eo` | non-EO allow-retry → exit **1**, still failed |
| `recovery/usage/missing-note` | recovery without `--note` → exit **2** |
| `recovery/usage/missing-id` | recovery without migration_id → exit **2** |
| `non-interactive/closed-stdin` | Closed stdin + `-h` under 2s, exit 0 |

## How to Run

```sh
cd /Users/xhd2015/Projects/xhd2015/mysql-migrate
# optional: export MIGRATE_MYSQL_DSN='user:pass@tcp(host:port)/db?multiStatements=true&...'
doctest vet ./tests/cli
doctest test ./tests/cli
```

MySQL must be reachable for status/plan/apply/recovery DB leaves (default
`localhost:9306` / `lifespan_db`). Offline leaves (help, unknown, usage,
recovery/usage, empty-note, closed-stdin) do not need MySQL. DB leaves **skip**
when the DSN is not reachable.

Classic TDD: `cli` is stub-only (`cli/doc.go`) until implementer lands `Run`.
Leaves must fail (compile or assertion RED) until:

```text
cli.Run(cfg migrate.Config, args []string) int
```

### CLI surface (locked)

```text
cli.Run(cfg, []string{"status"})
cli.Run(cfg, []string{"plan"})
cli.Run(cfg, []string{"apply"})                          # optional --to <id>
cli.Run(cfg, []string{"mark-done",   id, "--note", "..."})
cli.Run(cfg, []string{"mark-failed", id, "--note", "..."})
cli.Run(cfg, []string{"note",        id, "--note", "..."})
cli.Run(cfg, []string{"allow-retry", id, "--note", "..."})
```

No `--local` / `--remote`. Prefer **cfg only** for DSN and MigrationsDir
(tests pass both on Config; do not rely on `MIGRATE_MIGRATIONS_DIR` env).

### Exit codes

| Case | Exit |
|------|------|
| Help / success status, plan, apply, or recovery | **0** |
| Status/plan with `HasBlock`; apply refuse/exec fail; recovery biz error | **1** |
| Usage, unknown, missing DSN/MigrationsDir/id/note | **2** |

### Apply output contract

**Progress** (stdout, one line per attempted apply):

- success: contains `apply`, the `migration_id`, and `ok`
- failure: contains id (or `apply`) and `failed`

**Summary** (stdout): tokens `applied`, `failed`, `pending` with counts.

**Refuse-block**: stderr has **Error** mentioning **blocked**.

### Recovery contract

- All recovery commands require `--note` non-empty after trim.
- Missing migration_id or note → exit **2**.
- Happy path: exit **0**, log row reflects operation.
- `allow-retry` non-EO: exit **1**, stderr **Error**, row stays failed.
- Multi-step leaves use `FollowUpArgs` after primary exit **0**.

```go
import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/xhd2015/mysql-migrate/cli"
	"github.com/xhd2015/mysql-migrate/migrate"
)

// Request drives one CLI library invocation via cli.Run(cfg, args).
type Request struct {
	// Args is passed to cli.Run (without program name), e.g. []string{"-h"}.
	Args []string

	// CloseStdin when true installs a closed pipe as os.Stdin before Run.
	// When false, stdin is /dev/null (also non-blocking).
	CloseStdin bool

	// Config fields for migrate.Config (passed to cli.Run).
	// Prefer cfg only — MigrationsDir is NOT injected via env by Run.
	DSN           string
	MigrationsDir string
	ProgramName   string
	AppliedBy     string

	// FixtureIDs are migration_id values written by leaf Setup (basename without
	// .sql). Assert uses them to find rows in stdout / log; order is inventory
	// order (filename ascending).
	FixtureIDs []string

	// TableNames are optional MySQL table names created by apply fixtures.
	TableNames []string

	// SuccessDurationMS is the duration_ms seeded on a success log row when
	// the leaf wants status to surface that duration (0 = not asserted).
	SuccessDurationMS int

	// RecoveryNote is the operator --note string used by recovery happy-path
	// leaves (also present in Args). Asserts compare log.Note to this field.
	RecoveryNote string

	// FollowUpArgs when non-empty is a second cli.Run after primary exit 0,
	// with the same Config. Used for mark-done→status and allow-retry→apply.
	FollowUpArgs []string
}

// Response is the captured process-like outcome of cli.Run.
type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
	// Duration is wall time spent inside cli.Run (for hang detection).
	Duration time.Duration

	// Follow-up invocation (FollowUpExitCode=-1 when not run).
	FollowUpStdout   string
	FollowUpStderr   string
	FollowUpExitCode int // -1 if FollowUpArgs empty or primary exit != 0
	FollowUpDuration time.Duration
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	if req.Args == nil {
		req.Args = []string{}
	}

	cfg := buildConfig(req)
	stdout, stderr, code, dur, err := captureCLIRun(t, cfg, req.Args, req.CloseStdin)
	if err != nil {
		return nil, err
	}
	resp := &Response{
		Stdout:           stdout,
		Stderr:           stderr,
		ExitCode:         code,
		Duration:         dur,
		FollowUpExitCode: -1,
	}

	// Chain follow-up only after primary success (recovery → status/apply).
	if len(req.FollowUpArgs) > 0 && code == 0 {
		fo, fe, fc, fd, ferr := captureCLIRun(t, cfg, req.FollowUpArgs, req.CloseStdin)
		if ferr != nil {
			return nil, ferr
		}
		resp.FollowUpStdout = fo
		resp.FollowUpStderr = fe
		resp.FollowUpExitCode = fc
		resp.FollowUpDuration = fd
	}
	return resp, nil
}

func buildConfig(req *Request) migrate.Config {
	return migrate.Config{
		DSN:           req.DSN,
		MigrationsDir: req.MigrationsDir,
		ProgramName:   req.ProgramName,
		AppliedBy:     req.AppliedBy,
	}
}

// captureCLIRun redirects os.Stdout/os.Stderr/os.Stdin, calls cli.Run, restores.
func captureCLIRun(t *testing.T, cfg migrate.Config, args []string, closeStdin bool) (stdout, stderr string, exitCode int, dur time.Duration, err error) {
	t.Helper()

	rOut, wOut, err := os.Pipe()
	if err != nil {
		return "", "", 0, 0, err
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		_ = rOut.Close()
		_ = wOut.Close()
		return "", "", 0, 0, err
	}

	oldOut, oldErr, oldIn := os.Stdout, os.Stderr, os.Stdin
	os.Stdout, os.Stderr = wOut, wErr

	var stdinToClose *os.File
	if closeStdin {
		rIn, wIn, e := os.Pipe()
		if e != nil {
			os.Stdout, os.Stderr, os.Stdin = oldOut, oldErr, oldIn
			_ = rOut.Close()
			_ = wOut.Close()
			_ = rErr.Close()
			_ = wErr.Close()
			return "", "", 0, 0, e
		}
		_ = wIn.Close() // closed write end → immediate EOF on read
		os.Stdin = rIn
		stdinToClose = rIn
	} else {
		devNull, e := os.Open(os.DevNull)
		if e != nil {
			os.Stdout, os.Stderr, os.Stdin = oldOut, oldErr, oldIn
			_ = rOut.Close()
			_ = wOut.Close()
			_ = rErr.Close()
			_ = wErr.Close()
			return "", "", 0, 0, e
		}
		os.Stdin = devNull
		stdinToClose = devNull
	}

	outCh := make(chan string, 1)
	errCh := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, rOut)
		outCh <- buf.String()
	}()
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, rErr)
		errCh <- buf.String()
	}()

	start := time.Now()
	exitCode = cli.Run(cfg, args)
	dur = time.Since(start)

	_ = wOut.Close()
	_ = wErr.Close()
	stdout = <-outCh
	stderr = <-errCh
	_ = rOut.Close()
	_ = rErr.Close()
	if stdinToClose != nil {
		_ = stdinToClose.Close()
	}

	os.Stdout, os.Stderr, os.Stdin = oldOut, oldErr, oldIn
	return stdout, stderr, exitCode, dur, nil
}
```
