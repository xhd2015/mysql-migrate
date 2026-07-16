# mysql-migrate — plan status machine (P3)

Pure library tests for the migration **plan engine**: given an ordered list of
migration files (from inventory types) and log rows, compute per-migration
action (`skip` / `apply` / `blocked` / `deferred`) and `HasBlock` with **no I/O**.

Standalone doctest root under `tests/plan/` so inventory
stays independent while this package is missing.

Target package (implementer provides):

```text
github.com/xhd2015/mysql-migrate/migrate/plan
```

Uses inventory types only:

```text
github.com/xhd2015/mysql-migrate/migrate/inventory.MigrationFile
```

# DSN (Domain Specific Notion)

The **migration plan engine** is a pure function over **migration files** (from
inventory metadata) and **log rows** (prior apply attempts). For each file in
sorted file order it **classifies** an action: **apply** when there is no log,
empty status, or `"pending"`; **skip** when status is `success` and the log
hash matches the file hash **or** the log hash is empty; **blocked** when
status is `failed` or `unknown`, when status is stale `running` (normalized to
unknown), when EXACTLY-ONCE has any non-success terminal outcome, or when
success has a **hash mismatch** (never silent re-apply of the same id). It then
runs a **stop-chain**: after the first **blocked** item, later items that would
**apply** become **deferred**; **skip** items remain **skip**. The plan sets
`HasBlock` when any item is **blocked**. No MySQL, no filesystem, no log writes
inside `plan.Build`. File `ExactlyOnce` is the source of truth for EO rules.

## Version

0.0.2

## Decision Tree

```
tests/plan/                    [Request{Files, Logs}]
│                                            Run: plan.Build
├── empty/
│   └── no-files/                            # empty inputs → empty plan, no block
├── clear/                                   # HasBlock=false (skip/apply only)
│   ├── all-pending/                         # no logs → all apply
│   ├── all-success-matching-hash/           # success+hash match → all skip
│   ├── success-then-pending/                # skip then apply
│   ├── success-empty-log-hash-skips/        # success + empty log hash → skip
│   └── status-pending-or-empty-applies/     # status "" / "pending" → apply
└── blocked/                                 # HasBlock=true + deferred stop-chain
    ├── eo-failed-defers-later/              # EXACTLY-ONCE failed → block later
    ├── eo-unknown-blocks/                   # EXACTLY-ONCE unknown → blocked
    ├── non-eo-failed-defers-later/          # non-EO failed still blocks chain
    ├── unknown-status-blocks/               # status unknown → blocked
    ├── stale-running-as-unknown/            # running → effective unknown → blocked
    ├── success-hash-mismatch/               # success + hash differ → blocked
    └── skip-after-block-stays-skip/         # skip not rewritten to deferred
```

**Significance order:** chain outcome (empty | clear | blocked) → classification
rule / concrete file+log mix.

## Test Index

| Leaf | Description |
|------|-------------|
| `empty/no-files` | Empty files+logs → empty Items, HasBlock=false |
| `clear/all-pending` | Two files, no logs → apply, apply; HasBlock=false |
| `clear/all-success-matching-hash` | Success logs with matching hashes → all skip |
| `clear/success-then-pending` | First success skip, second no log apply |
| `clear/success-empty-log-hash-skips` | Success + empty log ContentSHA256 → skip |
| `clear/status-pending-or-empty-applies` | Log status `""` / `"pending"` → apply |
| `blocked/eo-failed-defers-later` | EO failed blocked; later pending deferred; HasBlock |
| `blocked/eo-unknown-blocks` | EO unknown → blocked + later deferred |
| `blocked/non-eo-failed-defers-later` | Non-EO failed blocks; later deferred |
| `blocked/unknown-status-blocks` | Non-EO unknown → blocked |
| `blocked/stale-running-as-unknown` | `running` → LogStatus unknown → blocked (EO + non-EO) |
| `blocked/success-hash-mismatch` | Success + hash mismatch → blocked + HashMismatch |
| `blocked/skip-after-block-stays-skip` | After blocked, later success remains skip |

## How to Run

```sh
cd /Users/xhd2015/Projects/xhd2015/mysql-migrate
doctest vet ./tests/plan
doctest test ./tests/plan
```

Classic RED: package is stub-only (`migrate/plan/doc.go`) until implementer
ports `LogRow`, `PlanItem`, `Plan`, `ItemAction`, and `Build`. Do not implement
here — designer owns tests only.

Public API expected by these tests:

```go
package plan

type LogRow struct {
    MigrationID   string
    Status        string // running | success | failed | unknown | pending | ""
    ExactlyOnce   bool   // snapshot; file.ExactlyOnce is source of truth for rules
    ContentSHA256 string
    DurationMS    int
    ErrorMessage  string
    Note          string
}

type ItemAction string // "skip" | "apply" | "blocked" | "deferred"

type PlanItem struct {
    MigrationID  string
    ExactlyOnce  bool
    Action       ItemAction
    Reason       string
    HashMismatch bool
    LogStatus    string // effective after stale-running → unknown
}

type Plan struct {
    Items    []PlanItem
    HasBlock bool
}

func Build(files []inventory.MigrationFile, logs []LogRow) Plan
```

```go
import (
	"fmt"
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

// Request carries in-memory files + log rows for plan.Build (no FS).
type Request struct {
	// Files is the migration inventory snapshot (order may be unsorted;
	// Build must emit Items in FileName ascending order).
	Files []inventory.MigrationFile

	// Logs are prior apply log rows; keyed by MigrationID inside Build.
	// Missing MigrationID = no row for that file.
	Logs []plan.LogRow
}

// PlanItemView is a plain snapshot of plan.PlanItem for assertions.
type PlanItemView struct {
	MigrationID  string
	ExactlyOnce  bool
	Action       string // "skip" | "apply" | "blocked" | "deferred"
	Reason       string
	HashMismatch bool
	LogStatus    string
}

// Response holds the pure plan outcome.
type Response struct {
	Items    []PlanItemView
	HasBlock bool
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	p := plan.Build(req.Files, req.Logs)
	items := make([]PlanItemView, 0, len(p.Items))
	for _, it := range p.Items {
		items = append(items, PlanItemView{
			MigrationID:  it.MigrationID,
			ExactlyOnce:  it.ExactlyOnce,
			Action:       string(it.Action),
			Reason:       it.Reason,
			HashMismatch: it.HashMismatch,
			LogStatus:    it.LogStatus,
		})
	}
	return &Response{Items: items, HasBlock: p.HasBlock}, nil
}
```
