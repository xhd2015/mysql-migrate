# mysql-migrate — log table + repository (P4)

Library tests for the tool-owned **`t_sql_migration_log`** table and the
**logrepo** package: ensure schema (idempotent), upsert lifecycle rows
(running → success/failed), list/get, and human recovery field updates
(mark-done, mark-failed-manual, set-note, allow-retry).

Standalone doctest root under `tests/logrepo/` so inventory / plan /
scaffold stay independent while this package is missing or RED.

Target package (implementer provides):

```text
github.com/xhd2015/mysql-migrate/migrate/logrepo
```

Live MySQL via `database/sql` + `github.com/go-sql-driver/mysql`. DSN from
env `MIGRATE_MYSQL_DSN` or the default local dev DSN (see SETUP).

# DSN (Domain Specific Notion)

The **migration log repository** owns a single MySQL table
**`t_sql_migration_log`** that records every apply attempt by **migration_id**.
A **caller** first **ensures** the table exists (`CREATE TABLE IF NOT EXISTS`) —
safe to call repeatedly. On apply, the runner **marks running** (upsert by
unique `migration_id`, storing exactly-once flag, content hash, and who
applied), then **marks success** (duration, finished time) or **marks failed**
(duration + error_message). Operators may **mark done** (force success with a
required note), **mark failed manually** (failed + note), **set note** without
changing status, or **allow retry** on an **exactly-once** failed row
(status → `pending` + note so plan can re-apply). Non-EO rows must not clear
via allow-retry. **List** returns all rows; **Get** returns one by migration_id.
Statuses: `running` | `success` | `failed` | `unknown` | `pending`.

## Version

0.0.2

## Decision Tree

```
tests/logrepo/                               [Request{Op, MigrationID, …}]
│                                            Run: logrepo via sql.Open(DSN)
├── ensure/
│   └── twice-idempotent/                    # EnsureTable ×2 OK; table exists
├── lifecycle/                               # MarkRunning → terminal status
│   ├── running-then-success/                # success + duration + hash + applied_by
│   ├── running-then-failed/                 # failed + error_message
│   └── unique-upsert/                       # second MarkRunning same id upserts
├── list/
│   └── after-inserts/                       # List contains seeded migration_ids
└── recovery/                                # human recovery / note ops (P8 CLI later)
    ├── mark-done/
    │   ├── with-note/                       # status=success + note
    │   └── empty-note-errors/               # empty note → error
    ├── mark-failed-manual/
    │   └── with-note/                       # status=failed + note
    ├── set-note-only/                       # note changes; status unchanged
    └── allow-retry/
        ├── exactly-once/                    # EO failed → pending + note
        └── non-eo-errors/                   # not ExactlyOnce → error
```

**Significance order:** operation class (ensure | lifecycle | list | recovery) →
outcome / validity (success vs failed, note present vs empty, EO vs non-EO) →
concrete field values.

## Test Index

| Leaf | Description |
|------|-------------|
| `ensure/twice-idempotent` | `EnsureTable` twice succeeds; `t_sql_migration_log` exists |
| `lifecycle/running-then-success` | MarkRunning then MarkSuccess → status success, duration, hash, applied_by |
| `lifecycle/running-then-failed` | MarkRunning then MarkFailed → failed + error_message |
| `lifecycle/unique-upsert` | Second MarkRunning same id upserts (one row); updated fields visible |
| `list/after-inserts` | List after multiple MarkRunning rows includes all migration_ids |
| `recovery/mark-done/with-note` | MarkDone with note → success + note |
| `recovery/mark-done/empty-note-errors` | MarkDone empty note → error |
| `recovery/mark-failed-manual/with-note` | MarkFailedManual → failed + note |
| `recovery/set-note-only` | SetNote changes note only; status stays prior value |
| `recovery/allow-retry/exactly-once` | EO failed → AllowRetry → pending + note |
| `recovery/allow-retry/non-eo-errors` | AllowRetry on non-EO row → error |

## How to Run

```sh
cd /Users/xhd2015/Projects/xhd2015/mysql-migrate
# optional: export MIGRATE_MYSQL_DSN='user:pass@tcp(host:port)/db?...'
doctest vet ./tests/logrepo
doctest test ./tests/logrepo
```

MySQL must be reachable at the resolved DSN (default local
`localhost:9306` / `lifespan_db`). Leaves **skip** when the DSN is not
reachable so pure-unit trees stay usable offline.

Classic TDD: `migrate/logrepo` is stub-only until implementer ports the API.
Leaves must fail (compile or assertion RED) until implementer lands:

```text
migrate/logrepo
```

Public API expected by these tests:

```go
package logrepo

import "database/sql"

// EnsureTable creates t_sql_migration_log if missing (IF NOT EXISTS). Idempotent.
func EnsureTable(db *sql.DB) error

// Row maps one log table row (align fields with plan.LogRow for later conversion).
type Row struct {
    MigrationID   string
    Status        string // running | success | failed | unknown | pending
    ExactlyOnce   bool
    ContentSHA256 string
    DurationMS    int
    ErrorMessage  string
    Note          string
    AppliedBy     string
    // optional timestamps may exist on the table; not required on Row for asserts
}

func List(db *sql.DB) ([]Row, error)
func Get(db *sql.DB, migrationID string) (Row, bool, error)

// Upsert lifecycle (unique on migration_id — prefer UPSERT / ON DUPLICATE KEY).
func MarkRunning(db *sql.DB, migrationID string, exactlyOnce bool, contentSHA256 string, appliedBy string) error
func MarkSuccess(db *sql.DB, migrationID string, durationMS int) error
func MarkFailed(db *sql.DB, migrationID string, durationMS int, errMsg string) error

// Human recovery (P8 CLI-wraps these; persistence required now).
// MarkDone / MarkFailedManual / AllowRetry require non-empty note.
func MarkDone(db *sql.DB, migrationID string, note string) error           // status=success + note
func MarkFailedManual(db *sql.DB, migrationID string, note string) error // status=failed + note
func SetNote(db *sql.DB, migrationID string, note string) error
func AllowRetry(db *sql.DB, migrationID string, note string) error // EO only; status→pending + note
```

**Table** `t_sql_migration_log` (tool-owned, not a numbered migration file):

| Column | Type |
|--------|------|
| id | BIGINT PK AI |
| migration_id | VARCHAR(255) NOT NULL UNIQUE |
| status | VARCHAR(32) NOT NULL |
| exactly_once | TINYINT(1) NOT NULL DEFAULT 0 |
| content_sha256 | VARCHAR(64) NOT NULL DEFAULT '' |
| started_at | DATETIME NULL |
| finished_at | DATETIME NULL |
| duration_ms | INT NOT NULL DEFAULT 0 |
| error_message | TEXT |
| note | TEXT |
| applied_by | VARCHAR(255) NOT NULL DEFAULT '' |
| create_time | DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP |
| update_time | DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP |

**plan.LogRow compatibility** (existing plan package): MigrationID, Status,
ExactlyOnce, ContentSHA256, DurationMS, ErrorMessage, Note. logrepo.Row must
expose these for later conversion; AppliedBy is logrepo-only.

**Isolation:** each leaf uses a unique `migration_id` prefix (`p4-` + session +
leaf slug). Leaves DELETE their id(s) before seed. Do not TRUNCATE the shared
dev table globally.

**Out of scope for P4:** CLI subcommands, applying migration SQL files,
plan.Build changes, inventory.

```go
import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/xhd2015/mysql-migrate/migrate/logrepo"
)

// defaultLocalDSN is the lifelog/local-dev MySQL DSN used when
// MIGRATE_MYSQL_DSN is unset. multiStatements not required for logrepo.
const defaultLocalDSN = "lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True"

// Request drives one logrepo scenario against MySQL.
type Request struct {
	// Op is the scenario dispatch key:
	// ensure_twice | lifecycle_success | lifecycle_failed | unique_upsert |
	// list | mark_done | mark_failed_manual | set_note | allow_retry
	Op string

	// MigrationID is the primary row key. Leaves set a unique isolated id.
	MigrationID string

	// ExactlyOnce stored by MarkRunning / seed.
	ExactlyOnce bool

	// ContentSHA256 stored by MarkRunning.
	ContentSHA256 string

	// AppliedBy stored by MarkRunning.
	AppliedBy string

	// DurationMS for MarkSuccess / MarkFailed.
	DurationMS int

	// ErrorMessage for MarkFailed.
	ErrorMessage string

	// Note for MarkDone / MarkFailedManual / SetNote / AllowRetry.
	Note string

	// Second* used by unique_upsert for the second MarkRunning call.
	SecondExactlyOnce   bool
	SecondContentSHA256 string
	SecondAppliedBy     string

	// ExtraMigrationIDs are additional rows seeded for list.
	ExtraMigrationIDs []string

	// SeedStatus for recovery leaves that need a pre-existing terminal row
	// before the recovery op: "running" | "success" | "failed" (via lifecycle helpers).
	// Empty means: seed with MarkRunning only (status running) before recovery op.
	SeedStatus string
}

// RowView is a plain snapshot of logrepo.Row for assertions.
type RowView struct {
	MigrationID   string
	Status        string
	ExactlyOnce   bool
	ContentSHA256 string
	DurationMS    int
	ErrorMessage  string
	Note          string
	AppliedBy     string
}

// Response holds table/row outcomes after the scenario.
type Response struct {
	// TableExists is true when information_schema lists t_sql_migration_log.
	TableExists bool

	// EnsureCallsOK is true when both EnsureTable calls returned nil (ensure_twice).
	EnsureCallsOK bool

	// Row is Get(MigrationID) after the op sequence; nil if not found / not fetched.
	Row *RowView

	// Found is whether Get found MigrationID.
	Found bool

	// Rows is the full List() result (callers filter by prefix in Assert).
	Rows []RowView

	// UniqueCount is COUNT(*) WHERE migration_id = MigrationID (unique_upsert).
	UniqueCount int
}

// resolveDSN returns MIGRATE_MYSQL_DSN if set, else defaultLocalDSN.
func resolveDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("MIGRATE_MYSQL_DSN")); dsn != "" {
		return dsn
	}
	return defaultLocalDSN
}

// openDB opens MySQL with the resolved DSN and pings.
func openDB(t *testing.T) (*sql.DB, error) {
	t.Helper()
	dsn := resolveDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping MySQL (DSN host from MIGRATE_MYSQL_DSN or default localhost:9306): %w", err)
	}
	return db, nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	if req.Op == "" {
		return nil, fmt.Errorf("empty Op")
	}

	db, err := openDB(t)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	switch req.Op {
	case "ensure_twice":
		if err := logrepo.EnsureTable(db); err != nil {
			return &Response{}, err
		}
		if err := logrepo.EnsureTable(db); err != nil {
			return &Response{}, err
		}
		exists, qerr := tableExists(db, "t_sql_migration_log")
		if qerr != nil {
			return &Response{EnsureCallsOK: true}, qerr
		}
		return &Response{EnsureCallsOK: true, TableExists: exists}, nil

	case "lifecycle_success":
		if err := ensureReady(db, req.MigrationID); err != nil {
			return nil, err
		}
		if err := logrepo.MarkRunning(db, req.MigrationID, req.ExactlyOnce, req.ContentSHA256, req.AppliedBy); err != nil {
			return &Response{}, err
		}
		if err := logrepo.MarkSuccess(db, req.MigrationID, req.DurationMS); err != nil {
			return &Response{}, err
		}
		return fetchRowResp(db, req.MigrationID)

	case "lifecycle_failed":
		if err := ensureReady(db, req.MigrationID); err != nil {
			return nil, err
		}
		if err := logrepo.MarkRunning(db, req.MigrationID, req.ExactlyOnce, req.ContentSHA256, req.AppliedBy); err != nil {
			return &Response{}, err
		}
		if err := logrepo.MarkFailed(db, req.MigrationID, req.DurationMS, req.ErrorMessage); err != nil {
			return &Response{}, err
		}
		return fetchRowResp(db, req.MigrationID)

	case "unique_upsert":
		if err := ensureReady(db, req.MigrationID); err != nil {
			return nil, err
		}
		if err := logrepo.MarkRunning(db, req.MigrationID, req.ExactlyOnce, req.ContentSHA256, req.AppliedBy); err != nil {
			return &Response{}, err
		}
		// Second MarkRunning same id must upsert, not fail with duplicate key.
		if err := logrepo.MarkRunning(db, req.MigrationID, req.SecondExactlyOnce, req.SecondContentSHA256, req.SecondAppliedBy); err != nil {
			return &Response{}, err
		}
		resp, err := fetchRowResp(db, req.MigrationID)
		if err != nil {
			return resp, err
		}
		n, cerr := countMigrationID(db, req.MigrationID)
		if cerr != nil {
			return resp, cerr
		}
		resp.UniqueCount = n
		return resp, nil

	case "list":
		if err := logrepo.EnsureTable(db); err != nil {
			return nil, err
		}
		// Seed primary + extras.
		ids := append([]string{req.MigrationID}, req.ExtraMigrationIDs...)
		for _, id := range ids {
			if id == "" {
				continue
			}
			_ = deleteMigrationID(db, id)
			if err := logrepo.MarkRunning(db, id, false, "hash-"+id, "list-seed"); err != nil {
				return &Response{}, err
			}
		}
		rows, err := logrepo.List(db)
		if err != nil {
			return &Response{}, err
		}
		return &Response{Rows: toRowViews(rows)}, nil

	case "mark_done":
		if err := seedThen(db, req); err != nil {
			return nil, err
		}
		if err := logrepo.MarkDone(db, req.MigrationID, req.Note); err != nil {
			return &Response{}, err
		}
		return fetchRowResp(db, req.MigrationID)

	case "mark_failed_manual":
		if err := seedThen(db, req); err != nil {
			return nil, err
		}
		if err := logrepo.MarkFailedManual(db, req.MigrationID, req.Note); err != nil {
			return &Response{}, err
		}
		return fetchRowResp(db, req.MigrationID)

	case "set_note":
		if err := seedThen(db, req); err != nil {
			return nil, err
		}
		if err := logrepo.SetNote(db, req.MigrationID, req.Note); err != nil {
			return &Response{}, err
		}
		return fetchRowResp(db, req.MigrationID)

	case "allow_retry":
		if err := seedThen(db, req); err != nil {
			return nil, err
		}
		if err := logrepo.AllowRetry(db, req.MigrationID, req.Note); err != nil {
			return &Response{}, err
		}
		return fetchRowResp(db, req.MigrationID)

	default:
		return nil, fmt.Errorf("unknown op %q", req.Op)
	}
}

func ensureReady(db *sql.DB, migrationID string) error {
	if err := logrepo.EnsureTable(db); err != nil {
		return err
	}
	if migrationID != "" {
		_ = deleteMigrationID(db, migrationID)
	}
	return nil
}

// seedThen ensures table, deletes id, seeds a row for recovery ops.
// SeedStatus: "" or "running" → MarkRunning only;
// "success" → MarkRunning + MarkSuccess;
// "failed" → MarkRunning + MarkFailed.
func seedThen(db *sql.DB, req *Request) error {
	if err := ensureReady(db, req.MigrationID); err != nil {
		return err
	}
	if err := logrepo.MarkRunning(db, req.MigrationID, req.ExactlyOnce, req.ContentSHA256, req.AppliedBy); err != nil {
		return err
	}
	switch strings.ToLower(strings.TrimSpace(req.SeedStatus)) {
	case "", "running":
		return nil
	case "success":
		return logrepo.MarkSuccess(db, req.MigrationID, req.DurationMS)
	case "failed":
		msg := req.ErrorMessage
		if msg == "" {
			msg = "seed-failed"
		}
		return logrepo.MarkFailed(db, req.MigrationID, req.DurationMS, msg)
	default:
		return fmt.Errorf("unknown SeedStatus %q", req.SeedStatus)
	}
}

func fetchRowResp(db *sql.DB, migrationID string) (*Response, error) {
	row, ok, err := logrepo.Get(db, migrationID)
	if err != nil {
		return &Response{}, err
	}
	if !ok {
		return &Response{Found: false}, nil
	}
	v := toRowView(row)
	return &Response{Found: true, Row: &v}, nil
}

func toRowView(r logrepo.Row) RowView {
	return RowView{
		MigrationID:   r.MigrationID,
		Status:        r.Status,
		ExactlyOnce:   r.ExactlyOnce,
		ContentSHA256: r.ContentSHA256,
		DurationMS:    r.DurationMS,
		ErrorMessage:  r.ErrorMessage,
		Note:          r.Note,
		AppliedBy:     r.AppliedBy,
	}
}

func toRowViews(rows []logrepo.Row) []RowView {
	out := make([]RowView, 0, len(rows))
	for _, r := range rows {
		out = append(out, toRowView(r))
	}
	return out
}

func tableExists(db *sql.DB, name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var n int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = DATABASE() AND table_name = ?`, name).Scan(&n)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func countMigrationID(db *sql.DB, migrationID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var n int
	err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM t_sql_migration_log WHERE migration_id = ?`, migrationID).Scan(&n)
	return n, err
}

func deleteMigrationID(db *sql.DB, migrationID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx,
		`DELETE FROM t_sql_migration_log WHERE migration_id = ?`, migrationID)
	return err
}
```
