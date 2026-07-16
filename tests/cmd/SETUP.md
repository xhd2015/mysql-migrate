# Scenario

**Feature**: thin `cmd/mysql-migrate` binary — less-flags globals, help, `cli.Run` hand-off

```
# session-built binary; flags/env → Config → cli.Run
operator -> mysql-migrate [--dsn] [--dir] <cmd> [args]
  -> less-flags parse globals
  -> migrate.Config{DSN, MigrationsDir, ProgramName}
  -> cli.Run(cfg, remain) -> stdout/stderr + exit code

# root help (binary-owned)
empty args | -h -> Usage (commands + --dsn/--dir) -> exit 0
```

## Preconditions

- Module: `github.com/xhd2015/mysql-migrate` (repo root `go.mod`).
- Binary package: `./cmd/mysql-migrate` (empty stub until implementer).
- Library: `cli.Run` already implemented under `cli/` (P5); this tree only
  locks main wiring.
- Global flags (less-flags): `--dsn`, `--dir`; help `-h` / `--help`.
- Env fallbacks: `MIGRATE_MYSQL_DSN`, `MIGRATE_MYSQL_DIR` (optional; flag wins).
- Session cache: `$TMPDIR/mysql-migrate-cmd-doctest-<DOCTEST_SESSION_ID>/`
  holds built binary + lock/ready markers (shared across parallel leaves).
- Default `ClearMigrateEnv=true` so ambient migrate env does not leak into
  missing-DSN leaves.
- Module root from this DOCTEST root: `DOCTEST_ROOT/../..` (`tests/cmd` → repo).
- Apply leaf needs MySQL (default `localhost:9306` / `lifespan_db`); skips
  when unreachable. Harness does **not** start containers.
- Apply DSN should allow multi-statement SQL (`multiStatements=true`).

## Steps

1. Root Setup builds (or reuses) the session binary and sets isolation defaults.
2. Leaves set `Args` (and fixtures / AssertDSN for apply).
3. Root `Run` execs `Bin` with controlled env; Assert checks exit + tokens /
   log side effects.

## Context

- Classic TDD: empty `main` builds and runs with exit 0 / no help text →
  assertion RED until implementer lands less-flags + `cli.Run`.
- Prefer flock session build over per-leaf `go run` for parallel leaves.
- Do not re-test full CLI matrix here (`tests/cli/` already seals it).

```go
import (
	"context"
	"database/sql"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/xhd2015/mysql-migrate/migrate/logrepo"
)

// defaultLocalDSN is the lifelog/local-dev MySQL DSN used when MIGRATE_MYSQL_DSN
// is unset for harness-side asserts / apply leaf skip detection.
const defaultLocalDSN = "lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True&multiStatements=true"

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	// Isolate from developer shell: strip migrate env unless a leaf opts out.
	// (Leaves that need env injection set ClearMigrateEnv and Env themselves;
	// default true so missing-DSN tests stay offline-stable.)
	req.ClearMigrateEnv = true
	req.Bin = buildBinaryOnce(t)
	return nil
}

// harnessDSN returns MIGRATE_MYSQL_DSN or the default local DSN.
// Used by harness (ensureMySQL / Assert), not automatically passed to the binary.
func harnessDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("MIGRATE_MYSQL_DSN")); dsn != "" {
		return dsn
	}
	return defaultLocalDSN
}

// sessionSlugPrefix returns a short kebab-safe prefix for migration slugs.
func sessionSlugPrefix() string {
	var b strings.Builder
	b.WriteString("p6")
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

func fixtureSlug(leaf, part string) string {
	return sessionSlugPrefix() + "-" + leaf + "-" + part
}

func fixtureTable(leaf, part string) string {
	return "t_cmd_" + sessionSlugPrefix() + "_" + leaf + "_" + part
}

func simpleFileName(seq int, slug string) string {
	return fmt.Sprintf("2026-07-16-%02d-%s.sql", seq, slug)
}

func migrationIDFromFile(fileName string) string {
	return strings.TrimSuffix(fileName, ".sql")
}

func writeMigration(t *testing.T, dir, fileName, body string) string {
	t.Helper()
	path := filepath.Join(dir, fileName)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write migration %s: %v", path, err)
	}
	return migrationIDFromFile(fileName)
}

func createTableSQL(tableName string) string {
	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s` (\n  id INT NOT NULL PRIMARY KEY\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;\n",
		tableName,
	)
}

func contentSHA256(body string) string {
	sum := sha256.Sum256([]byte(body))
	return hex.EncodeToString(sum[:])
}

// ensureMySQL pings harness DSN; skips the leaf when MySQL is unavailable.
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
		t.Skipf("MySQL not reachable at harness DSN (skip apply leaf): %v", pingErr)
	}
}

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

func deleteLogIDs(t *testing.T, db *sql.DB, ids ...string) {
	t.Helper()
	for _, id := range ids {
		if id == "" {
			continue
		}
		_, _ = db.Exec(`DELETE FROM t_sql_migration_log WHERE migration_id = ?`, id)
	}
}

func dropTables(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()
	for _, name := range tables {
		if name == "" {
			continue
		}
		_, _ = db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", name))
	}
}

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

func requireLogStatus(t *testing.T, db *sql.DB, migrationID, wantStatus string) {
	t.Helper()
	row, found, err := logrepo.Get(db, migrationID)
	if err != nil {
		t.Fatalf("logrepo.Get %q: %v", migrationID, err)
	}
	if !found {
		t.Fatalf("log: want status %q for %q, but no row", wantStatus, migrationID)
	}
	if row.Status != wantStatus {
		t.Fatalf("log: %q status got %q want %q", migrationID, row.Status, wantStatus)
	}
}

func requireExit(t *testing.T, resp *Response, want int) {
	t.Helper()
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != want {
		t.Fatalf("exit: got %d want %d\nstdout=%q\nstderr=%q",
			resp.ExitCode, want, resp.Stdout, resp.Stderr)
	}
}

// requireRootHelpTokens checks Usage + command list + global flags.
func requireRootHelpTokens(t *testing.T, stdout string) {
	t.Helper()
	if !strings.Contains(stdout, "Usage") {
		t.Fatalf("stdout must contain Usage:\n%s", stdout)
	}
	for _, cmd := range []string{
		"status",
		"plan",
		"apply",
		"mark-done",
		"mark-failed",
		"note",
		"allow-retry",
	} {
		if !strings.Contains(stdout, cmd) {
			t.Fatalf("root help must list %q:\n%s", cmd, stdout)
		}
	}
	// Global flags owned by the binary (less-flags).
	if !strings.Contains(stdout, "--dsn") {
		t.Fatalf("root help must mention --dsn:\n%s", stdout)
	}
	if !strings.Contains(stdout, "--dir") {
		t.Fatalf("root help must mention --dir:\n%s", stdout)
	}
	if strings.Contains(stdout, "--local") || strings.Contains(stdout, "--remote") {
		t.Fatalf("root help must not mention --local/--remote:\n%s", stdout)
	}
}

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
