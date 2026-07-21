# Scenario

**Feature**: thin `cmd/mysql-migrate` binary — less-flags globals, edge open/Wrap, `cli.Run` hand-off

```
# session-built binary; flags/env → edge open → Config.DB → cli.Run
operator -> mysql-migrate [--dsn] [--dir] <cmd> [args]
  -> less-flags parse globals
  -> if dsn: sql.Open → Ping → sqlexec.Wrap → cfg.DB
  -> migrate.Config{DB, MigrationsDir, ProgramName}  # no DSN field
  -> cli.Run(cfg, remain) -> stdout/stderr + exit code

# root help (binary-owned)
empty args | -h -> Usage (commands + --dsn/--dir) -> exit 0
```

## Preconditions

- Module: `github.com/xhd2015/mysql-migrate` (repo root `go.mod`).
- Binary package: `./cmd/mysql-migrate` — edge opens DSN only; core Config is DB-only.
- Library: `cli.Run` never `sql.Open`; requires non-nil `cfg.DB` for DB cmds.
- Global flags (less-flags): `--dsn`, `--dir`; help `-h` / `--help`.
- Env fallbacks: `MIGRATE_MYSQL_DSN`, `MIGRATE_MYSQL_DIR` (optional; flag wins).
- Process-local binary/cache via in-memory mutex (one-process suite; not in-memory mutex)
  holds built binary + lock/ready markers (shared across leaves in one process).
- MySQL-touching leaves acquire in-process `mysqlExclusiveMu` so
  `status/ensure-created` can drop `t_sql_migration_log` without racing
  sibling DB leaves in this tree.
- Default `ClearMigrateEnv=true` so ambient migrate env does not leak into
  missing-DSN leaves.
- Module root from this DOCTEST root: `d.DOCTEST_ROOT/../..` (`tests/cmd` → repo).
- DB leaves need MySQL (default `localhost:9306` / `lifespan_db`); skip when
  unreachable. Harness does **not** start containers.
- Apply DSN should allow multi-statement SQL (`multiStatements=true`).

## Steps

1. Root Setup builds (or reuses) the session binary and sets isolation defaults.
2. Leaves set `Args` (and fixtures / AssertDSN for DB leaves).
3. Root `Run` execs `Bin` with controlled env; Assert checks exit + tokens /
   log side effects.

## Context

- P2 edge contract: `--dsn`/`--dir` → Open → Wrap → `cli.Run`; no DSN on Config.
- Prefer process-local in-memory once build over per-leaf `go run` for parallel leaves.
- Do not re-test full CLI matrix here (`tests/cli/` already seals it).

```go
import (
	"sync"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/xhd2015/mysql-migrate/migrate/logrepo"
	"github.com/xhd2015/doctest/session"
)

// defaultLocalDSN is the lifelog/local-dev MySQL DSN used when MIGRATE_MYSQL_DSN
// is unset for harness-side asserts / apply leaf skip detection.
const defaultLocalDSN = "lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True&multiStatements=true"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	// Isolate from developer shell: strip migrate env unless a leaf opts out.
	// (Leaves that need env injection set ClearMigrateEnv and Env themselves;
	// default true so missing-DSN tests stay offline-stable.)
	req.ClearMigrateEnv = true
	req.Bin = buildBinaryOnce(t, d)
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
func sessionSlugPrefix(d *session.Doctest) string {
	var b strings.Builder
	b.WriteString("p6")
	for _, r := range d.DOCTEST_SESSION_ID {
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

func fixtureSlug(d *session.Doctest, leaf, part string) string {
	return sessionSlugPrefix(d) + "-" + leaf + "-" + part
}

func fixtureTable(d *session.Doctest, leaf, part string) string {
	return "t_cmd_" + sessionSlugPrefix(d) + "_" + leaf + "_" + part
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

// Process-local MySQL reachability memo (one-process; not session flock).
var (
	mysqlExclusiveMu sync.Mutex
	ensureMySQLMu  sync.Mutex
	ensureMySQLDid bool
	ensureMySQLErr error
)

func ensureMySQL(t *testing.T, d *session.Doctest) {
	t.Helper()
	ensureMySQLMu.Lock()
	defer ensureMySQLMu.Unlock()
	if ensureMySQLDid {
		if ensureMySQLErr != nil {
			t.Skipf("MySQL not reachable at harness DSN (skip apply leaf): %v", ensureMySQLErr)
		}
		return
	}
	ensureMySQLDid = true
	dsn := harnessDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		ensureMySQLErr = err
		t.Skipf("MySQL not reachable at harness DSN (skip apply leaf): %v", ensureMySQLErr)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ensureMySQLErr = db.PingContext(ctx)
	cancel()
	_ = db.Close()
	if ensureMySQLErr != nil {
		t.Skipf("MySQL not reachable at harness DSN (skip apply leaf): %v", ensureMySQLErr)
	}
}

func openLocalDB(t *testing.T, d *session.Doctest) *sql.DB {
	t.Helper()
	ensureMySQL(t, d)
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

// acquireMySQLExclusive serialises MySQL-touching leaves in-process.
// Call from every leaf Setup that touches the shared harness MySQL.
func acquireMySQLExclusive(t *testing.T) {
	t.Helper()
	mysqlExclusiveMu.Lock()
	t.Cleanup(func() { mysqlExclusiveMu.Unlock() })
}


// dropMigrationLogTable drops t_sql_migration_log on the harness DB.
// Caller must hold acquireMySQLExclusive so sibling leaves do not race.
func dropMigrationLogTable(t *testing.T, d *session.Doctest) {
	t.Helper()
	db := openLocalDB(t, d)
	defer db.Close()
	if _, err := db.Exec(`DROP TABLE IF EXISTS t_sql_migration_log`); err != nil {
		t.Fatalf("DROP t_sql_migration_log: %v", err)
	}
}

// requireMissingDSNUsage checks exit 2 and Error tokens for nil-DB path.
// Binary leaves cfg.DB nil when no flag/env DSN; cli prints missing DB.
func requireMissingDSNUsage(t *testing.T, resp *Response) {
	t.Helper()
	requireExit(t, resp, 2)
	combined := resp.Stdout + "\n" + resp.Stderr
	lower := strings.ToLower(combined)
	if !strings.Contains(lower, "dsn") &&
		!strings.Contains(lower, "missing") &&
		!strings.Contains(lower, "db") {
		t.Fatalf("missing DSN usage must mention dsn, missing, or db:\nstdout=%q\nstderr=%q",
			resp.Stdout, resp.Stderr)
	}
	if strings.Contains(combined, "--local") || strings.Contains(combined, "--remote") {
		t.Fatalf("usage error must not require --local/--remote:\n%s", combined)
	}
}
```
