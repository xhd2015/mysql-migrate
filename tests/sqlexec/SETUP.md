# Scenario

**Feature**: context-first `sqlexec.DB` facade + Config.DB-only engine contract

```
# caller owns *sql.DB; library only sees sqlexec.DB
caller -> sql.Open(DSN) -> sqlexec.Wrap(sqlDB) -> DB

# execute with context (never DSN inside migrate library)
caller -> DB.Exec|Query|QueryRow(ctx, sql, args...) -> Result|Rows|Row
caller -> DB.Close() -> further ops error

# Config carries DB, not DSN
migrate.Config{DB, MigrationsDir, ProgramName, AppliedBy}
cli.Run(cfg{DB:nil}, ["status"]) -> exit 2 usage (missing DB)
```

## Preconditions

- Module: `github.com/xhd2015/mysql-migrate` (repo root `go.mod`).
- Target package (to implement):
  `migrate/sqlexec`
  import `github.com/xhd2015/mysql-migrate/migrate/sqlexec`
- Config contract: `migrate.Config` has `DB sqlexec.DB` and **no** `DSN` field.
- Live leaves use `database/sql` + `github.com/go-sql-driver/mysql` only in the
  **harness** to obtain `*sql.DB` for `Wrap`. Production packages under test
  must not open DSN themselves.
- **DSN resolution** (harness only):
  1. Env `MIGRATE_MYSQL_DSN` if non-empty
  2. Else default:
     `lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True`
- Offline leaves: `interface/methods-present`, `config/*` (no MySQL).
- Live leaves: root Setup pings MySQL; **skips** when unreachable.
- Module root from this DOCTEST root: `d.DOCTEST_ROOT/..`.
- Isolation: live leaves use unique table names via `tableName(leaf)` so parallel
  leaves do not collide. Prefer `DROP TABLE IF EXISTS` cleanup.
- Session cache may record MySQL readiness once per `doctest test` run (flock + ready).

## Steps

1. Root Setup defaults nothing heavy; live-branch Setup calls `ensureMySQL`.
2. Leaf Setup sets `req.Op`, optional `Table` / `SeedValues` / `MigrationsDir`.
3. Root `Run` dispatches on `Op` against `sqlexec` (and Config/cli for config leaves).
4. Leaf Assert checks interface shape, result rows, errors, or CLI exit codes.

## Context

- Classic TDD: importing `migrate/sqlexec` fails compile until implementer lands
  the package (RED).
- `cli` already exists with DSN-based Config; `config/nil-db-cli-usage` goes RED
  until Config.DB-only + cli require DB.
- Do not modify lifelog. Do not write production package code in design phase.
- Parallel-safe: per-leaf tables; session-scoped MySQL ready marker.

```go
import (
	"sync"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xhd2015/doctest/session"
)

// defaultLocalDSNPing matches DOCTEST.md defaultLocalDSN for harness connectivity.
const defaultLocalDSNPing = "lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Root: normalize slices; branches set Op and whether MySQL is required.
	// Offline leaves (interface, config) must not require MySQL at root.
	if req.SeedValues == nil {
		req.SeedValues = []int64{}
	}
	t.Logf("sqlexec root setup: session=%s repo=%s", d.DOCTEST_SESSION_ID, repoRoot(t, d))
	return nil
}

func repoRoot(t *testing.T, d *session.Doctest) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join(d.DOCTEST_ROOT, ".."))
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

// Process-local MySQL reachability memo (one-process; not session flock).
var (
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

// tableName builds an isolated InnoDB table name for a live leaf.
// MySQL identifier-safe: prefix t_sqx_ + short session + leaf slug.
func tableName(d *session.Doctest, leaf string) string {
	sid := d.DOCTEST_SESSION_ID
	var b strings.Builder
	b.WriteString("t_sqx_")
	n := 0
	for _, r := range sid {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			n++
		} else if r >= 'A' && r <= 'Z' {
			b.WriteRune(r - 'A' + 'a')
			n++
		}
		if n >= 8 {
			break
		}
	}
	if n == 0 {
		b.WriteString("x")
	}
	b.WriteByte('_')
	// leaf slug: keep alnum only
	for _, r := range leaf {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if r >= 'A' && r <= 'Z' {
			b.WriteRune(r - 'A' + 'a')
		} else if r == '-' || r == '_' {
			b.WriteByte('_')
		}
	}
	name := b.String()
	// MySQL max identifier 64; keep headroom.
	if len(name) > 60 {
		name = name[:60]
	}
	return name
}

// requireNoHarnessErr fails when Run itself failed (as opposed to OpErr in Response).
func requireNoHarnessErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
}
```
