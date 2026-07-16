# Scenario

**Feature**: non-interactive migrate CLI library — help, routing, status/plan, apply, recovery via `cli.Run(cfg, args)`

```
# operator/main invokes library with Config + argv; CLI dispatches without hanging
operator -> cli.Run(cfg, args) -> stdout/stderr + exit code (never os.Exit)

# help / usage (offline)
args -h | <cmd> -h -> Usage with ProgramName (exit 0)
unknown cmd | nil cfg.DB for DB ops -> exit 2

# status / plan (local DB + fixture migrations dir on Config)
cfg{DSN, MigrationsDir} + status|plan
  -> sql.Open(DSN) -> logrepo.EnsureTable -> inventory.ListDir
  -> logrepo.List -> plan.Build -> table on stdout
  -> exit 0 if !HasBlock else 1; hash mismatch => stderr warning:

# apply (mutating)
cfg + apply [--to id]
  -> same plan pipeline
  -> if blocked: refuse, Error blocked, exit 1
  -> else for each apply: MarkRunning -> Exec SQL -> MarkSuccess|MarkFailed
  -> progress + summary on stdout; exit 0|1

# recovery (human gate — no migration SQL; no --local/--remote)
mark-done|mark-failed|note|allow-retry <id> --note "..."
  -> logrepo.MarkDone|MarkFailedManual|SetNote|AllowRetry
  -> exit 0; allow-retry non-EO => exit 1 Error
  -> missing note|id => exit 2
```

## Preconditions

- Module: `github.com/xhd2015/mysql-migrate` (repo root `go.mod`).
- Target package: `cli` import `github.com/xhd2015/mysql-migrate/cli`.
- Config type: `migrate.Config` from `github.com/xhd2015/mysql-migrate/migrate`.
- Supporting packages: `inventory`, `plan`, `logrepo` (already present).
- Public API: `func Run(cfg migrate.Config, args []string) int` (args **without**
  program name; **never** `os.Exit`).
- Exit codes: `0` success/help/clear plan/successful apply/recovery, `1` HasBlock
  or apply failure/refuse or recovery biz error, `2` usage / unknown / missing
  DSN or MigrationsDir / missing id or note.
- Full command set (root help): `status`, `plan`, `apply`, `mark-done`,
  `mark-failed`, `note`, `allow-retry`.
- **No** `--local` / `--remote` flags anywhere.
- Config: tests set harness `DSN` (open string only) + `MigrationsDir` on
  Request for DB leaves. Root `Run` does `sql.Open` + `sqlexec.Wrap` into
  `cfg.DB`. Empty harness DSN → nil `cfg.DB` (usage leaves). Prefer cfg only;
  harness does not set `MIGRATE_MIGRATIONS_DIR`.
- Module root from this DOCTEST root: `DOCTEST_ROOT/..`.
- Status/plan/apply/recovery DB leaves need MySQL (default `localhost:9306` /
  `lifespan_db`). Harness **skips** when DSN not reachable (no podman start).
- Isolation: fixture migration filenames embed a session-scoped slug prefix so
  seeded `t_sql_migration_log` rows do not collide across parallel leaves.
  Prefer DELETE-by-id cleanup over TRUNCATE. Apply fixtures use unique table
  names (`t_mig_<session>_<leaf>_<part>`) and DROP TABLE on cleanup.
- Recovery multi-step leaves may set `FollowUpArgs` (status or apply) after
  primary exit 0; `RecoveryNote` holds the operator note for log asserts.
- Apply DSN must allow multi-statement SQL (`multiStatements=true`).

## Steps

1. Root Setup defaults ProgramName/AppliedBy; leaves set Args, harness DSN
   (open string), MigrationsDir, fixtures, seeds, FollowUpArgs.
2. Root `Run` builds `migrate.Config` with `cfg.DB = sqlexec.Wrap(...)` when
   harness DSN is set, redirects stdio, calls `cli.Run`, restores, returns
   captured text + exit code + duration; on primary exit 0 with
   `FollowUpArgs`, runs a second `cli.Run` with the same Config.
3. Leaf Assert checks exit code, stdout/stderr tokens, and for apply/recovery
   leaves: log status/note + optional table existence side effects.

## Context

- Classic TDD: importing `cli.Run` fails compile until implementer lands the API
  beyond the P1 stub package (`cli/doc.go`).
- Session cache may record MySQL readiness once per `doctest test` run (flock + ready).
- Parallel-safe: each leaf captures/restores process stdio.

```go
import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/xhd2015/mysql-migrate/migrate/logrepo"
)

// defaultLocalDSN is the lifelog/local-dev MySQL DSN used when MIGRATE_MYSQL_DSN
// is unset. multiStatements=true is required for apply multi-statement files.
const defaultLocalDSN = "lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True&multiStatements=true"

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	// Defaults for help text / applied_by identity; leaves may override.
	if req.ProgramName == "" {
		req.ProgramName = "mysql-migrate"
	}
	if req.AppliedBy == "" {
		req.AppliedBy = "cli-doctest"
	}
	// CloseStdin defaults false; non-interactive leaf sets true.
	// DSN / MigrationsDir / FixtureIDs / TableNames / seeds left for DB leaves.
	return nil
}

func repoRoot(t *testing.T) string {
	t.Helper()
	// Root at tests/cli → module root is ..
	root, err := filepath.Abs(filepath.Join(DOCTEST_ROOT, ".."))
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	return root
}

// harnessDSN returns MIGRATE_MYSQL_DSN or the default local DSN (with multiStatements).
func harnessDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("MIGRATE_MYSQL_DSN")); dsn != "" {
		return dsn
	}
	return defaultLocalDSN
}

// fillConfigForDB ensures MySQL is up and fills harness req.DSN when empty.
// Root Run wraps that DSN into cfg.DB via sqlexec.Wrap (Config has no DSN field).
// Leaves still set MigrationsDir themselves.
func fillConfigForDB(t *testing.T, req *Request) {
	t.Helper()
	ensureMySQL(t)
	if strings.TrimSpace(req.DSN) == "" {
		req.DSN = harnessDSN() // harness-only open string → Wrap in buildConfig
	}
	if req.ProgramName == "" {
		req.ProgramName = "mysql-migrate"
	}
	if req.AppliedBy == "" {
		req.AppliedBy = "cli-doctest"
	}
}

// sessionSlugPrefix returns a short kebab-safe prefix for migration slugs.
func sessionSlugPrefix() string {
	var b strings.Builder
	b.WriteString("p5")
	for _, r := range DOCTEST_SESSION_ID {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r - 'A' + 'a')
		}
		if b.Len() >= 10 {
			break
		}
	}
	if b.Len() < 4 {
		b.WriteString("x")
	}
	return b.String()
}

// fixtureSlug builds a unique kebab slug: <sessionPrefix>-<leaf>-<part>
func fixtureSlug(leaf, part string) string {
	return sessionSlugPrefix() + "-" + leaf + "-" + part
}

// fixtureTable builds a unique MySQL table name for apply side-effect asserts.
func fixtureTable(leaf, part string) string {
	return "t_mig_" + sessionSlugPrefix() + "_" + leaf + "_" + part
}

// simpleFileName builds YYYY-MM-DD-NN-<slug>.sql
func simpleFileName(seq int, slug string) string {
	return fmt.Sprintf("2026-07-16-%02d-%s.sql", seq, slug)
}

// eoFileName builds YYYY-MM-DD-NN-[EXACTLY-ONCE]-<slug>.sql
func eoFileName(seq int, slug string) string {
	return fmt.Sprintf("2026-07-16-%02d-[EXACTLY-ONCE]-%s.sql", seq, slug)
}

// migrationIDFromFile returns basename without .sql
func migrationIDFromFile(fileName string) string {
	return strings.TrimSuffix(fileName, ".sql")
}

// writeMigration writes body to dir/fileName and returns migration_id.
func writeMigration(t *testing.T, dir, fileName, body string) string {
	t.Helper()
	path := filepath.Join(dir, fileName)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write migration %s: %v", path, err)
	}
	return migrationIDFromFile(fileName)
}

// createTableSQL returns idempotent CREATE TABLE IF NOT EXISTS for name.
func createTableSQL(tableName string) string {
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s` (\n  id INT NOT NULL PRIMARY KEY\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n",
		tableName,
	)
}

// contentSHA256 returns lowercase hex SHA-256 of body bytes.
func contentSHA256(body string) string {
	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
}

func sessionCacheDir() string {
	return filepath.Join(os.TempDir(), "mysql-migrate-cli-doctest-"+DOCTEST_SESSION_ID)
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
// so pure-unit trees remain usable offline. Does not start containers.
func ensureMySQL(t *testing.T) {
	t.Helper()
	cache := sessionCacheDir()
	lock := filepath.Join(cache, "mysql.lock")
	ready := filepath.Join(cache, "mysql.ready")
	dsn := harnessDSN()

	var pingErr error
	err := withFileLock(t, lock, func() error {
		if _, statErr := os.Stat(ready); statErr == nil {
			db, err := sql.Open("mysql", dsn)
			if err == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				pingErr = db.PingContext(ctx)
				cancel()
				db.Close()
				if pingErr == nil {
					return nil
				}
				// stale ready marker — fall through to re-ping path below
				_ = os.Remove(ready)
			}
		}
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			pingErr = err
			return nil
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		pingErr = db.PingContext(ctx)
		cancel()
		db.Close()
		if pingErr == nil {
			return os.WriteFile(ready, []byte("ok"), 0o644)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("ensureMySQL lock: %v", err)
	}
	if pingErr != nil {
		t.Skipf("MySQL not reachable at harness DSN (skip DB leaf): %v", pingErr)
	}
}

// openLocalDB opens sql.DB with harnessDSN after ensureMySQL; caller must Close.
func openLocalDB(t *testing.T) *sql.DB {
	t.Helper()
	ensureMySQL(t)
	db, err := sql.Open("mysql", harnessDSN())
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		t.Fatalf("ping MySQL: %v", err)
	}
	return db
}

// deleteLogIDs removes seeded rows for isolation (best-effort).
func deleteLogIDs(t *testing.T, db *sql.DB, ids ...string) {
	t.Helper()
	for _, id := range ids {
		if id == "" {
			continue
		}
		_, _ = db.Exec(`DELETE FROM t_sql_migration_log WHERE migration_id = ?`, id)
	}
}

// dropTables drops fixture tables (best-effort) for apply leaf isolation.
func dropTables(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()
	for _, name := range tables {
		if name == "" {
			continue
		}
		_, _ = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", name))
	}
}

// tableExists reports whether name exists in the current database.
func tableExists(t *testing.T, db *sql.DB, name string) bool {
	t.Helper()
	var n int
	err := db.QueryRow(
		`SELECT COUNT(*) FROM information_schema.tables
		 WHERE table_schema = DATABASE() AND table_name = ?`,
		name,
	).Scan(&n)
	if err != nil {
		t.Fatalf("tableExists %q: %v", name, err)
	}
	return n > 0
}

// logStatus returns status string and whether the log row exists.
func logStatus(t *testing.T, db *sql.DB, migrationID string) (status string, ok bool) {
	t.Helper()
	row, found, err := logrepo.Get(db, migrationID)
	if err != nil {
		t.Fatalf("logrepo.Get %q: %v", migrationID, err)
	}
	if !found {
		return "", false
	}
	return row.Status, true
}

// requireLogStatus fails unless migrationID has wantStatus in the log.
func requireLogStatus(t *testing.T, db *sql.DB, migrationID, wantStatus string) {
	t.Helper()
	got, ok := logStatus(t, db, migrationID)
	if !ok {
		t.Fatalf("log: want status %q for %q, but no row", wantStatus, migrationID)
	}
	if got != wantStatus {
		t.Fatalf("log: %q status got %q want %q", migrationID, got, wantStatus)
	}
}

// requireLogNote fails unless migrationID exists and Note equals want.
func requireLogNote(t *testing.T, db *sql.DB, migrationID, want string) {
	t.Helper()
	row, found, err := logrepo.Get(db, migrationID)
	if err != nil {
		t.Fatalf("logrepo.Get %q: %v", migrationID, err)
	}
	if !found {
		t.Fatalf("log: want note for %q, but no row", migrationID)
	}
	if row.Note != want {
		t.Fatalf("log note %q: got %q want %q", migrationID, row.Note, want)
	}
}

// requireFollowUpExit fails unless a follow-up ran with want exit code.
func requireFollowUpExit(t *testing.T, resp *Response, want int) {
	t.Helper()
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.FollowUpExitCode < 0 {
		t.Fatalf("follow-up did not run (primary exit=%d stdout=%q stderr=%q)",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if resp.FollowUpExitCode != want {
		t.Fatalf("follow-up exit: got %d want %d\nstdout=%q\nstderr=%q",
			resp.FollowUpExitCode, want, resp.FollowUpStdout, resp.FollowUpStderr)
	}
}

// seedSuccess upserts running→success with hash and duration; optional note via SetNote.
func seedSuccess(t *testing.T, db *sql.DB, migrationID string, exactlyOnce bool, hash string, durationMS int, note string) {
	t.Helper()
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	deleteLogIDs(t, db, migrationID)
	if err := logrepo.MarkRunning(db, migrationID, exactlyOnce, hash, "p5-cli-seed"); err != nil {
		t.Fatalf("MarkRunning %s: %v", migrationID, err)
	}
	if err := logrepo.MarkSuccess(db, migrationID, durationMS); err != nil {
		t.Fatalf("MarkSuccess %s: %v", migrationID, err)
	}
	if strings.TrimSpace(note) != "" {
		if err := logrepo.SetNote(db, migrationID, note); err != nil {
			t.Fatalf("SetNote %s: %v", migrationID, err)
		}
	}
	t.Cleanup(func() { deleteLogIDs(t, db, migrationID) })
}

// seedFailed upserts running→failed with hash and error message.
func seedFailed(t *testing.T, db *sql.DB, migrationID string, exactlyOnce bool, hash string, durationMS int, errMsg string) {
	t.Helper()
	if _, err := logrepo.EnsureTable(db); err != nil {
		t.Fatalf("EnsureTable: %v", err)
	}
	deleteLogIDs(t, db, migrationID)
	if err := logrepo.MarkRunning(db, migrationID, exactlyOnce, hash, "p5-cli-seed"); err != nil {
		t.Fatalf("MarkRunning %s: %v", migrationID, err)
	}
	if err := logrepo.MarkFailed(db, migrationID, durationMS, errMsg); err != nil {
		t.Fatalf("MarkFailed %s: %v", migrationID, err)
	}
	t.Cleanup(func() { deleteLogIDs(t, db, migrationID) })
}

// stdoutHasActionNearID checks that stdout mentions migrationID and action.
func stdoutHasActionNearID(stdout, migrationID, action string) bool {
	if !strings.Contains(stdout, migrationID) {
		return false
	}
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, migrationID) && strings.Contains(line, action) {
			return true
		}
	}
	return strings.Contains(stdout, action)
}

// requireExit fails unless resp.ExitCode matches want.
func requireExit(t *testing.T, resp *Response, want int) {
	t.Helper()
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != want {
		t.Fatalf("exit: got %d want %d\nstdout=%q\nstderr=%q", resp.ExitCode, want, resp.Stdout, resp.Stderr)
	}
}

// requireNotStub fails if stdout/stderr still say "not implemented".
func requireNotStub(t *testing.T, resp *Response, cmd string) {
	t.Helper()
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if strings.Contains(combined, "not implemented") {
		t.Fatalf("%s still stubbed:\nstdout=%q\nstderr=%q", cmd, resp.Stdout, resp.Stderr)
	}
}

// stdoutHasApplyProgress looks for a progress line mentioning id and token.
func stdoutHasApplyProgress(stdout, migrationID, token string) bool {
	if !strings.Contains(stdout, migrationID) {
		return false
	}
	for _, line := range strings.Split(stdout, "\n") {
		if !strings.Contains(line, migrationID) {
			continue
		}
		if strings.Contains(line, token) {
			return true
		}
	}
	return strings.Contains(stdout, token)
}
```
