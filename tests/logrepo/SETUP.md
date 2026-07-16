# Scenario

**Feature**: tool-owned `t_sql_migration_log` + logrepo read/write lifecycle

```
# ensure schema once (idempotent CREATE TABLE IF NOT EXISTS)
caller -> logrepo.EnsureTable(db) -> t_sql_migration_log

# apply lifecycle: running -> success | failed (unique migration_id upsert)
caller -> MarkRunning(id, eo, hash, by) -> row status=running
caller -> MarkSuccess(id, duration) | MarkFailed(id, duration, err) -> terminal

# recovery for operators (notes required where specified)
caller -> MarkDone / MarkFailedManual / SetNote / AllowRetry(EO only)

# read back
caller -> Get(id) | List() -> rows for plan/status
```

## Preconditions

- Module: `github.com/xhd2015/mysql-migrate` (repo root `go.mod`).
- Target package (to implement):
  `migrate/logrepo`
  import `github.com/xhd2015/mysql-migrate/migrate/logrepo`
- Live DB via `database/sql` + `github.com/go-sql-driver/mysql` (no `target` package).
- **DSN resolution** (document for operators / CI):
  1. Env `MIGRATE_MYSQL_DSN` if non-empty
  2. Else default:
     `lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True`
- Root Setup pings MySQL once; **skips** the leaf when unreachable (MySQL
  optional for offline pure-unit suites).
- Module root from this DOCTEST root: `DOCTEST_ROOT/..`.
- Isolation: leaf `MigrationID` values must use prefix from `idPrefix()` so
  parallel leaves do not collide. Prefer DELETE by id over TRUNCATE.
- Session cache may record MySQL readiness once per `doctest test` run (flock + ready).

## Steps

1. Root Setup pings resolved DSN (`ensureMySQL`); skip leaf if down.
2. Leaf Setup sets `req.Op`, unique `MigrationID`, and op-specific fields.
3. `Run` opens DB, ensures/seeds as needed, runs the logrepo API, returns
   row snapshot(s).
4. Leaf Assert checks status, fields, uniqueness, or errors.

## Context

- Classic TDD: importing `migrate/logrepo` fails compile until implementer lands
  the full API beyond the P1 stub (RED).
- All leaves need MySQL when the package exists; skip if DSN not reachable.
- plan.LogRow field subset must remain representable on logrepo.Row for later P6.
- Do not modify inventory/plan/scaffold trees or production packages in design phase.

```go
import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// defaultLocalDSNPing matches DOCTEST.md defaultLocalDSN for harness connectivity.
const defaultLocalDSNPing = "lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True"

func Setup(t *testing.T, req *Request) error {
	// All logrepo leaves hit MySQL — skip cleanly when DSN is not reachable.
	ensureMySQL(t)
	return nil
}

func repoRoot(t *testing.T) string {
	t.Helper()
	// Root at tests/logrepo → module root is ..
	root, err := filepath.Abs(filepath.Join(DOCTEST_ROOT, ".."))
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	return root
}

// harnessDSN returns MIGRATE_MYSQL_DSN or the default local DSN.
func harnessDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("MIGRATE_MYSQL_DSN")); dsn != "" {
		return dsn
	}
	return defaultLocalDSNPing
}

// idPrefix returns a session-scoped migration_id prefix for isolation.
func idPrefix() string {
	// Keep short: migration_id VARCHAR(255); session ids can be long.
	sid := DOCTEST_SESSION_ID
	if len(sid) > 12 {
		sid = sid[:12]
	}
	return "p4-" + sid + "-"
}

// leafMigrationID builds a unique migration_id for a leaf slug.
func leafMigrationID(slug string) string {
	return idPrefix() + slug
}

func sessionCacheDir() string {
	return filepath.Join(os.TempDir(), "mysql-migrate-logrepo-doctest-"+DOCTEST_SESSION_ID)
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

// ensureMySQL pings the resolved DSN. Skips the leaf when MySQL is unavailable
// so pure-unit trees remain usable offline. Does not start containers (this
// standalone repo has no podman-compose harness).
func ensureMySQL(t *testing.T) {
	t.Helper()
	cache := sessionCacheDir()
	lock := filepath.Join(cache, "mysql.lock")
	ready := filepath.Join(cache, "mysql.ready")
	dsn := harnessDSN()

	var pingErr error
	err := withFileLock(t, lock, func() error {
		if _, statErr := os.Stat(ready); statErr == nil {
			// Re-verify quickly; ready may be stale across reboots.
			db, err := sql.Open("mysql", dsn)
			if err == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				pingErr = db.PingContext(ctx)
				cancel()
				db.Close()
				if pingErr == nil {
					return nil
				}
				_ = os.Remove(ready)
			}
		}
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			pingErr = err
			return nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		pingErr = db.PingContext(ctx)
		cancel()
		db.Close()
		if pingErr == nil {
			_ = os.WriteFile(ready, []byte("ok"), 0o644)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ensureMySQL lock: %v", err)
	}
	if pingErr != nil {
		t.Skipf("MySQL not reachable at resolved DSN (set MIGRATE_MYSQL_DSN or start local MySQL on :9306): %v", pingErr)
	}
}

// requireRow fails unless resp has a found row.
func requireRow(t *testing.T, resp *Response) *RowView {
	t.Helper()
	if resp == nil {
		t.Fatal("nil response")
	}
	if !resp.Found || resp.Row == nil {
		t.Fatal("expected Found row, got none")
	}
	return resp.Row
}

// assertStatus fails unless row status matches want.
func assertStatus(t *testing.T, row *RowView, want string) {
	t.Helper()
	if row.Status != want {
		t.Fatalf("status: got %q want %q", row.Status, want)
	}
}

// assertNonEmpty fails if s is empty (for notes / error messages).
func assertNonEmpty(t *testing.T, field, s string) {
	t.Helper()
	if s == "" {
		t.Fatalf("%s: expected non-empty", field)
	}
}

// mustErr fails if err is nil.
func mustErr(t *testing.T, err error, context string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error, got nil", context)
	}
	if err.Error() == "" {
		t.Fatalf("%s: expected non-empty error message", context)
	}
}
```
