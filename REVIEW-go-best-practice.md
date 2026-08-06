# Go Best-Practice Review — mysql-migrate

**Scope:** package layout, CLI design, flag handling, and alignment with
`go-best-practice` topics (`flags-parsing/*`, `cli/*`, `kool-create`,
`cmd-exec`, `go-embed-assets` as applicable).  
**Mode:** review only — no product fixes implemented.  
**Date:** 2026-08-06  
**Module:** `github.com/xhd2015/mysql-migrate` (~1.5k LOC Go)

---

## Executive summary

This is a well-shaped Go migration operator. The **DSN-free library core**
(`migrate.Config` + `cli.Run` + pure `inventory`/`plan` packages) is a strong
architectural choice: the binary is a thin process edge that opens MySQL and
injects `sqlexec.DB`. Subcommand surface, exit codes (0/1/2), non-interactive
recovery, and doctest coverage are mature for the size of the project.

The main gaps against go-best-practice are concentrated in **flag parsing
consistency** and a few **CLI UX** issues:

| Area | Verdict |
|------|---------|
| Package layout / DSN-free core | Strong |
| Binary globals + `StopOnFirstArg` | Good start (`flags-parsing/subcommand`) |
| Subcommand flags via less-flags | **Missing** — hand-rolled parsers in `cli` |
| Help at every level | Present, but **library-leaky** and **blocked by DSN open** |
| Apply streaming progress | Good (`cli/streaming`) |
| Apply dry-run gate | Partial (`plan` ≈ preview; no `--dry-run` on apply) |
| Color / `NO_COLOR` | Not present (optional for this tool) |
| `go:embed` assets | N/A (no generated UI); DDL is a string const |
| `kool-create` / `cmd-exec` | N/A to current product code |

**Highest-value fixes (recommended order):** (1) do not open/ping DSN for
help-only paths; (2) parse subcommand flags with less-flags instead of
hand scanners; (3) promote `less-flags` to a direct `go.mod` require;
(4) rewrite operator-facing help so it describes flags/env, not `cfg.DB`;
(5) optional `--dry-run` on apply as a side-effect gate on the existing
pipeline.

---

## Project snapshot

```text
cmd/mysql-migrate/     process edge: less-flags globals, sql.Open, cli.Run
cli/                   library CLI: dispatch, status/plan/apply/recovery
migrate/               Config (DB + MigrationsDir + ProgramName + AppliedBy)
migrate/sqlexec/       context-first DB facade + Wrap(*sql.DB)
migrate/inventory/     pure on-disk file grammar + hash (no DB)
migrate/plan/          pure plan status machine (no I/O)
migrate/logrepo/       t_sql_migration_log ensure/lifecycle/recovery
tests/                 doctest trees (cli, cmd, inventory, plan, logrepo, …)
```

**Dependencies:** `github.com/go-sql-driver/mysql` (direct),
`github.com/xhd2015/less-flags` (**used directly by main, marked `// indirect`**).

**Flag ownership today:**

| Layer | Flags | Parser |
|-------|--------|--------|
| Binary (`cmd/mysql-migrate`) | `--dsn`, `--dir`, `-h/--help` | less-flags + `StopOnFirstArg` + `HelpNoExit` |
| CLI library (`cli.Run`) | `apply --to`, recovery `--note`, subcommand `-h` | **hand-rolled** loops |
| Config (not flags) | `DB`, `MigrationsDir`, `ProgramName`, `AppliedBy` | injected |

---

## Findings (by severity)

### High

#### H1. Subcommand flags are hand-parsed instead of less-flags

**Topic:** `flags-parsing`, `flags-parsing/types`, `flags-parsing/subcommand`

**Where:** `cli/cli.go` — `parseOptionalTo`, `parseIDAndNote`, `hasUnknownFlags`,
`wantsHelp`, `nonFlagPositional`

**What:** The binary correctly uses less-flags for globals. The CLI library
reimplements flag scanning for `--to`, `--note`, help detection, and unknown
flags. That diverges from the recipe: **every subcommand handler should run
its own `lessflags.Parse`** over its own option set (with `Help("-h,--help", …)`
or `HelpFunc` + `HelpNoExit` when returning exit codes).

**Why it matters:**

- Drift risk: equals forms (`--to=id`, `--note=…`), unknown flags, and help
  behavior are maintained twice.
- Subcommands do not get less-flags error messages / consistent parse failures.
- Adding flags (timeout, dry-run, color, verbose) will grow custom scanners.

**Recommended change:**

```go
// apply — sketch aligned with flags-parsing/subcommand
func runApply(cfg migrate.Config, program string, args []string) int {
    var to string
    remain, err := lessflags.String("--to", &to).
        HelpFunc("-h,--help", func() { printApplyHelp(os.Stdout, program) }).
        HelpNoExit().
        Parse(args)
    if err != nil {
        if err == lessflags.ErrHelp {
            return ExitOK
        }
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        return ExitUsage
    }
    if len(remain) > 0 {
        // unexpected positionals
        return ExitUsage
    }
    // requireDB / requireMigrationsDir / doApply...
}
```

Same pattern for recovery:

```go
lessflags.String("--note", &note).
    HelpFunc("-h,--help", ...).
    HelpNoExit().
    Parse(args)
// remain[0] = migration_id; enforce len(remain)==1 and non-empty note
```

For flag-less `status` / `plan`, still use less-flags with only
`Help(...)` (or empty parse + help) so unknown flags are rejected by the
parser, not a custom `hasUnknownFlags`.

**Note on library purity:** less-flags is already a product dependency via
the binary. Pulling it into `cli` is consistent; alternatively keep `cli`
flag-free and parse all flags in `cmd` then pass typed options — but today
flags already live in `cli`, so less-flags belongs there.

---

#### H2. Help paths open/ping DSN when `--dsn` is present

**Topic:** `cli` (UX), `flags-parsing/subcommand` (help must always answer)

**Where:** `cmd/mysql-migrate/main.go` — open/ping runs for every remaining
args path after global parse, including subcommand help.

**Evidence (smoke):**

```text
$ mysql-migrate --dsn 'bad' status -h
Error: open DSN: invalid DSN: ...
# exit 1 — no Usage printed

$ mysql-migrate --dsn 'user:pass@tcp(127.0.0.1:1)/db' apply -h
Error: open DSN: dial tcp 127.0.0.1:1: connect: connection refused
# exit 1 — help unreachable
```

**Why it matters:** Operators and scripts run `<cmd> -h` / `--help` offline
and in broken-env debugging. Help must not depend on MySQL connectivity.
The subcommand recipe requires every level to answer help; the binary edge
currently gates that on a successful Ping.

**Recommended change:**

1. After global parse, if remain is empty or is pure help
   (`-h` / `--help` / `help`) → print root help, exit 0, **no open**.
2. If remain is `<subcommand> -h|--help` (and no other work args), either:
   - open nothing and dispatch help-only into `cli`, or
   - detect help in main before open and call a help-only path.
3. Only `sql.Open` + `Ping` when the selected command actually needs DB
   (or always open for non-help DB commands).

Minimal shape:

```go
if isHelpOnly(remain) {
    return cli.Run(migrate.Config{ProgramName: programName}, remain)
}
// else open DSN when needed...
```

Also avoid opening for recovery/status help when DSN is set but invalid.

---

### Medium

#### M1. Operator-facing help leaks library Config vocabulary

**Topic:** `cli` / `flags-parsing/subcommand` (help text quality)

**Where:** `cli/cli.go` `printRootUsage`, per-command help strings

**What:** When users run the **binary** as `mysql-migrate status -h` (or
`mysql-migrate help`), they see:

```text
Uses cfg.DB and cfg.MigrationsDir from Config (no target flags).
```

and root library usage:

```text
Config (passed by the caller, not CLI flags):
  DB            sqlexec.DB ...
  MigrationsDir  Directory of *.sql migration files
```

Operators of the binary care about `--dsn` / `--dir` / env vars, not
`cfg.DB`. Library embedders care about Config. One help surface serves both
poorly.

**Recommended change:**

- Binary root help (already good in `printRootHelp`) stays flag-centric.
- Either:
  - **A)** Teach `cli` about “operator mode” vs “library mode” via
    `cfg.ProgramName` / a `HelpStyle` field, or
  - **B)** Keep library help abstract, but have `cmd` own all operator help
    text and only call `cli` for non-help execution, or
  - **C)** Dual help: subcommand help mentions both:
    - Operator: `--dsn` / `--dir` (global, before command) or env
    - Library: inject `Config.DB` / `Config.MigrationsDir`

Prefer **C** or binary-owned help for anything shipped as `mysql-migrate`.
Library package docs (`cli/doc.go`, README Architecture) already explain
Config well.

---

#### M2. `less-flags` is a direct import but `// indirect` in go.mod

**Topic:** `flags-parsing` (module hygiene)

**Where:** `go.mod`

```go
require github.com/go-sql-driver/mysql v1.10.0

require (
    ...
    github.com/xhd2015/less-flags v1.0.2 // indirect
)
```

**What:** `cmd/mysql-migrate` imports `github.com/xhd2015/less-flags`
directly (`go list` / `go mod why` confirm), yet the require is marked
indirect. That confuses consumers and tools that treat indirect deps as
optional.

**Recommended change (trivial):**

```sh
go get github.com/xhd2015/less-flags@v1.0.2
# or: go mod edit -require=github.com/xhd2015/less-flags@v1.0.2 && go mod tidy
```

Ensure a top-level:

```go
require (
    github.com/go-sql-driver/mysql v1.10.0
    github.com/xhd2015/less-flags v1.0.2
)
```

---

#### M3. Global flags only work *before* the subcommand (footgun)

**Topic:** `flags-parsing/subcommand` (`StopOnFirstArg` contract)

**Where:** `cmd/mysql-migrate/main.go` — `StopOnFirstArg()`; README says
flags are “before the subcommand”.

**Evidence:**

```text
$ mysql-migrate status --dsn 'x'
Error: unexpected arguments for status: --dsn x

$ mysql-migrate plan --dir ./migrations
Error: unexpected arguments for plan: --dir ./migrations
```

**What:** This is **correct** for the chosen `StopOnFirstArg` design, and
matches common kubectl-like “globals then command” tools. It is still a
UX footgun for users who expect `mysql-migrate status --dsn …`.

**Recommended change (docs + help, not necessarily redesign):**

1. Root and subcommand help should state explicitly:

   ```text
   Global flags must appear before the subcommand:
     mysql-migrate --dsn … --dir … status
   Not: mysql-migrate status --dsn …
   ```

2. Optionally improve the error when a known global flag appears after the
   command: `Error: --dsn is a global flag; place it before "status"`.

3. Do **not** silently accept globals after the command unless you redesign
   flag ownership (collect globals both sides — more complex; see
   `flags-parsing/collect` only if you intentionally forward filtered argv).

---

#### M4. No apply `--dry-run` (plan is related but not the same)

**Topic:** `cli/dry-run`

**Where:** `cli` apply path; `plan` subcommand

**What:** Best practice: **one pipeline**, gate side effects — do not fork a
separate dry-run implementation. Today:

- `plan` / `status` share `plan.Build` with apply (good: same classification).
- `apply` always mutates (MarkRunning / Exec / MarkSuccess|Failed).
- There is no `apply --dry-run` that walks the apply queue and prints
  `[dry-run] would apply <id>` without writes.

`plan` answers “what would apply decide?” but does not exercise the apply
loop’s `--to` filtering, progress line format, or summary counts in the
same path.

**Recommended change:**

```go
// Single doApply(cfg, toID, dryRun bool)
// After plan.Build + HasBlock check:
for _, item := range applyQueue {
    if dryRun {
        fmt.Fprintf(os.Stdout, "[dry-run] would apply %s\n", item.MigrationID)
        // honor --to; no MarkRunning / Exec
        continue
    }
    // existing MarkRunning → Exec → Mark*
}
```

Keep **one** function; do not add `doApplyDryRun` that reimplements inventory
or plan. Document that `plan` is the human table view; `apply --dry-run` is
the gated apply pipeline.

---

#### M5. Duplicated plan-build orchestration in status/plan vs apply

**Topic:** `cli/dry-run` (one pipeline), maintainability

**Where:** `runStatusOrPlan` and `doApply` both:

1. `ensureMigrationLog`
2. `inventory.ListDir`
3. `logrepo.List` → convert to `plan.LogRow`
4. `plan.Build`

**Recommended change:** Extract something like:

```go
func loadPlan(cfg migrate.Config) (plan.Plan, []inventory.MigrationFile, map[string]logrepo.Row, int)
```

Then status/plan print; apply gates side effects. This is the same structural
move that makes `--dry-run` cheap and correct.

---

#### M6. Binary always Pings; help and offline recovery usage mixed with open policy

**Topic:** `cli` UX

**Related to H2.** Even for non-help commands, Ping-on-start is fine for an
operator tool. Document that a reachable DSN is required for all DB
subcommands. Consider delaying open until after subcommand dispatch so
usage errors (`mark-done` missing `--note`) can be reported without MySQL
when DSN is set-but-broken.

Today recovery usage parsing correctly happens **inside** `cli` *after* main
has already opened. Order today:

```text
parse globals → open+ping if dsn → cli.Run → parse note/id
```

Better:

```text
parse globals → dispatch → if needsDB { open+ping } → run
// or: cli parses usage first; main opens only when cli signals need
```

Library-first design makes “main opens always” simpler; at least skip open
for help (H2) and consider skipping for pure usage failures if easy.

---

### Low

#### L1. No color / `NO_COLOR` policy

**Topic:** `cli/color`

**What:** Errors and `warning:` lines are plain text. Fine for a migration
tool used in CI logs. Optional enhancement:

- `--color` / `--no-color` + auto TTY + `NO_COLOR`
- Red `Error:`, yellow `warning:`, green `ok` tokens
- Never color if machine-readable output is added later

Not required unless operator UX polish is a goal.

---

#### L2. Apply streaming is good; status/plan table is batched (OK)

**Topic:** `cli/streaming`

**What:** Apply already prints per-migration progress then a summary —
aligned with “stream as work proceeds.” Status/plan correctly buffer for
`tabwriter` alignment (recipe allows two-pass for tables).

No change required. If plans grow huge, consider fixed-width columns or
NDJSON (`--json`) later; full JSON arrays only for small fixed payloads.

---

#### L3. Hard-coded timeouts; no `--timeout` Duration flag

**Topic:** `flags-parsing/types` (`Duration`)

**Where:** apply Exec 5m; logrepo ops 30s; binary Ping 15s

**Recommended change (optional):** global or apply-level
`lessflags.Duration("--timeout", &d)` with defaults. Prefer `**time.Duration`
or document zero = default if you need unset detection.

---

#### L4. `AppliedBy` has no binary flag/env

**Topic:** CLI completeness

**Where:** `migrate.Config.AppliedBy`; binary never sets it → defaults to
`"mysql-migrate"`.

**Recommended change:** optional `--applied-by` / `MIGRATE_MYSQL_APPLIED_BY`
for multi-operator audit trails. Low priority.

---

#### L5. `cli` hard-codes `os.Stdout` / `os.Stderr`

**Topic:** testability / library UX

**What:** `Run` documents writing to os stdout/stderr; doctests swap
`os.Stdout`. Injecting `io.Writer` (or a small `cli.IO` struct) would make
unit tests cleaner without fd redirection. Not wrong; incremental.

---

#### L6. `logrepo` accepts `any` (`sqlexec.DB` or `*sql.DB`)

**Topic:** package boundary clarity

**Where:** `logrepo.asDB`

**What:** Public surface prefers `sqlexec.DB`; `*sql.DB` remains for harness
compat. Slightly softens the “never *sql.DB in library” story. Prefer
`sqlexec.DB` only in signatures long-term; keep wrap in tests/cmd.

---

#### L7. `mysql-migrate help` vs `mysql-migrate --help` differ

**What:**

- `--help` / empty argv → binary `printRootHelp` (flags + commands)
- bare `help` subcommand → `cli` root usage (Config-oriented)

Confusing dual root help. Treat `help` as alias of root operator help in
main before open, or make `cli`’s help command match binary text when
`ProgramName` is the operator binary.

---

#### L8. multiStatements not documented on binary help

**Where:** README / apply comment mention multi-statement DSN; root help does
not.

**Recommended change:** one line under `--dsn` in help and README:

```text
Prefer multiStatements=true on the DSN so multi-statement migration files run.
```

---

### Informational / not applicable

#### I1. `kool-create`

Scaffolding skill for **new** projects (`server`, `go-react`, …). This repo
is an established CLI module with a deliberate DSN-free layout. No need to
re-scaffold. If starting a sibling tool, `kool create server <name>` is fine
as a blank slate — then apply the patterns already proven here (thin
`cmd/`, library packages, less-flags + `StopOnFirstArg`).

#### I2. `cmd-exec` (`xgo/support/cmd`)

Product code does not shell out to external processes (correct for a
migration engine). Doctest harness uses `os/exec` to run the built binary —
appropriate for integration tests; no need to switch harness to
`xgo/support/cmd` unless you want Debug logging of test invocations.

#### I3. `go-embed-assets`

Recipe targets **generated UI/extension** trees and release hydrate. This
project has no such assets. `logrepo` DDL is a small `const` string — fine.
Optional future: `//go:embed ddl/t_sql_migration_log.sql` for keeping DDL
next to consumer migrations, with a committed file (not placeholder hydrate).
Do **not** force the browser-agent multi-layer asset pattern here.

#### I4. Package layout strengths (keep)

- **Edge vs core:** only `cmd/mysql-migrate` imports the MySQL driver and
  opens DSNs; library packages stay DSN-free. Excellent for embedders and
  tests.
- **Pure packages:** `inventory` and `plan` have no DB — easy to doctest.
- **`sqlexec` facade:** context-first API; Wrap at the edge.
- **Exit codes:** 0 success/help, 1 business, 2 usage — consistent and tested.
- **Non-interactive recovery:** required `--note`, never blocks on stdin.
- **Doctest trees:** seal CLI and binary contracts separately (`tests/cli`
  vs `tests/cmd`).

#### I5. `go test ./...` vs doctest

CI runs both. There are **no** `*_test.go` unit files in packages; coverage
is doctest-driven. That is a project convention, not a go-best-practice
violation. Optional later: small pure unit tests for `inventory.ParseFileName`
/ `plan.Build` without MySQL.

---

## Topic checklist

| Topic | Applicable? | Project status | Action |
|-------|-------------|----------------|--------|
| `flags-parsing` | Yes | Globals only | Use less-flags in every subcommand handler |
| `flags-parsing/subcommand` | Yes | Partial (`StopOnFirstArg` + manual help) | H1, H2, M1, M3, L7 |
| `flags-parsing/types` | Yes | Strings only | Optional Duration; use `*string`/`**T` as needed |
| `flags-parsing/cut` | No | No opaque trailing exec | — |
| `flags-parsing/collect` | Maybe | Only if forwarding filtered argv | Not needed now |
| `cli/dry-run` | Yes | `plan` shares Build; no apply gate | M4, M5 |
| `cli/streaming` | Yes | Apply streams; tables batch | Keep; L2 |
| `cli/color` | Optional | None | L1 |
| `cli/skill-cli` | No | Not a skill host | — |
| `cli/inline-tui-mouse` | No | Non-interactive CLI | — |
| `cmd-exec` | No (product) | No external commands | I2 |
| `go-embed-assets` | No | No generated assets | I3 |
| `kool-create` | No (existing repo) | — | I1 |

---

## Recommended change backlog (implementation order)

When implementation is authorized, prefer this order so each step is
verifiable against existing doctests:

1. **H2** — Skip DSN open/ping for help-only invocations (binary + `cmd`
   doctests: help leaves still exit 0; add a leaf with bad `--dsn` +
   `status -h` → Usage exit 0).
2. **M2** — `go mod tidy` / promote `less-flags` to direct require.
3. **H1** — Replace hand parsers with less-flags in `cli` (`--to`, `--note`,
   status/plan unknown-flag rejection). Keep exit codes and help strings
   stable for sealed tests; update only if help text is intentionally
   improved.
4. **M1 + L7 + L8** — Operator-facing help rewrite (flags/env, global flag
   order, multiStatements note); unify `help` vs `--help`.
5. **M3** — Document global-before-command; friendlier misplaced-flag errors.
6. **M5 → M4** — Extract shared `loadPlan`; add `apply --dry-run` as
   side-effect gate on the same loop.
7. **L3 / L4 / L1** — Timeout, AppliedBy, color as polish.

**Docs-only (safe anytime):** README already states global flags before
subcommand; expand with the misplaced-flag example and multiStatements;
mention exit codes table.

---

## Architecture notes (positive patterns to preserve)

```text
                    ┌─────────────────────────┐
  operator argv ──► │ cmd/mysql-migrate       │  less-flags globals
                    │  sql.Open → Ping → Wrap │  StopOnFirstArg
                    └───────────┬─────────────┘
                                │ migrate.Config{DB, MigrationsDir, …}
                                ▼
                    ┌─────────────────────────┐
                    │ cli.Run (no os.Exit,    │
                    │         no sql.Open)    │
                    └───────────┬─────────────┘
           ┌────────────────────┼────────────────────┐
           ▼                    ▼                    ▼
     inventory.ListDir    logrepo.*             plan.Build
     (pure disk)          (sqlexec.DB)          (pure)
```

This split matches good Go CLI practice even beyond the skill recipes:
testable pure core, injectable I/O at the edge, and a library entry that
never owns process lifetime (`os.Exit` only in `main`).

---

## Conclusion

mysql-migrate already follows several go-best-practice ideas well:
**subcommand dispatch with `StopOnFirstArg`**, **help at multiple levels**,
**streamed apply progress**, **env fallbacks with flag precedence**, and a
**clean package layout** with a DSN-free core.

The review’s actionable thrust is to **finish the less-flags story inside
`cli`**, **never block help on MySQL**, **fix module require hygiene**, and
optionally **gate apply side effects with `--dry-run`** on the same pipeline
as live apply. No redesign of inventory/plan/logrepo is required for best-
practice alignment.
)
