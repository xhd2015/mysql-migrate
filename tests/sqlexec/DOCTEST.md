# mysql-migrate — `migrate/sqlexec` + Config.DB-only (P1)

Library tests for the context-first SQL execution facade **`sqlexec`** and the
**Config.DB-only** engine contract: the migrate library never opens a DSN;
callers inject a live `sqlexec.DB` (typically via `sqlexec.Wrap(*sql.DB)`).

Standalone doctest root under `tests/sqlexec/` so inventory / plan / logrepo /
cli / scaffold trees stay independent while this package is missing (classic
RED).

Target packages (implementer provides):

```text
github.com/xhd2015/mysql-migrate/migrate/sqlexec
github.com/xhd2015/mysql-migrate/migrate          # Config.DB field; no DSN
```

Related consumer contracts (documented + exercised lightly here; full harness
updates live in `tests/cli` / `tests/logrepo` / `tests/scaffold`):

- `logrepo` functions take `sqlexec.DB` (not `*sql.DB`)
- `cli.Run` never `sql.Open`; requires non-nil `cfg.DB` for DB subcommands
- nil `cfg.DB` → usage error (exit **2**)

# DSN (Domain Specific Notion)

**Participants**

- **Caller** — tests, CLI main (P2), or host app that already owns a
  `*database/sql.DB` connection.
- **`sqlexec` package** — thin facade over SQL execution with **context on
  every call**. Exposes interface **`DB`** and constructor **`Wrap`**.
- **`DB`** — executable handle: `Exec`, `Query`, `QueryRow`, `Close`. Does
  **not** know about DSNs, drivers, or `sql.Open`.
- **`Result`** — outcome of `Exec` (`LastInsertId`, `RowsAffected`).
- **`Rows`** — multi-row cursor from `Query` (`Next`, `Scan`, `Close`, `Err`).
- **`Row`** — single-row accessor from `QueryRow` (`Scan`).
- **`Wrap`** — adapts a standard `*sql.DB` into `sqlexec.DB` by forwarding to
  `*Context` methods on `database/sql`.
- **`migrate.Config`** — engine config for library callers: **`DB`**
  (`sqlexec.DB`), `MigrationsDir`, `ProgramName`, `AppliedBy`. **No DSN
  field** — opening connections is outside the library.
- **CLI library** — `cli.Run(cfg, args)` uses `cfg.DB` only; never calls
  `sql.Open`. Missing DB on a DB subcommand is a **usage** failure.

**Behaviors**

- Caller obtains `*sql.DB` externally (env DSN, pool, test harness), then
  `db := sqlexec.Wrap(sqlDB)`.
- `db.Exec(ctx, query, args...)` runs a statement and returns `Result` or error.
- `db.Query(ctx, query, args...)` returns `Rows` for multi-row reads; empty
  result sets are success with `Next()==false`.
- `db.QueryRow(ctx, query, args...).Scan(...)` reads one row; zero rows surface
  as `sql.ErrNoRows` (via `errors.Is`).
- `db.Close()` releases the wrapped handle; further Exec/Query fail.
- Zero-value `migrate.Config{}` has nil `DB` and empty string identity fields.
- `cli.Run` with nil `cfg.DB` on `status` (and other DB subcommands) exits **2**
  with an Error mentioning DB / missing config — not a hang and not a biz exit.

## Version

0.0.2

## Decision Tree

Split first on **surface under test** (largest behavioral impact), then on
**outcome class** (success vs error / empty):

```
tests/sqlexec/                               [Request{Op, Table, SQL, …}]
│                                            Run: sqlexec (+ Config / cli for config leaves)
├── interface/
│   └── methods-present/                     # DB/Result/Rows/Row + Wrap compile surface
├── wrap/
│   └── returns-db/                          # Wrap(*sql.DB) non-nil; SELECT 1 works
├── exec/
│   ├── success/                             # CREATE/INSERT; RowsAffected ≥ 1
│   └── bad-sql/                             # invalid SQL → error
├── query/
│   ├── multi-row/                           # INSERT n; Query scans all rows
│   └── empty/                               # Query no rows → Next false, Err nil
├── query-row/
│   ├── one-row/                             # QueryRow Scan one value OK
│   └── no-rows/                             # QueryRow Scan → errors.Is(..., sql.ErrNoRows)
├── close/
│   └── after-close-exec-errors/             # Close then Exec fails
└── config/
    ├── db-field-no-dsn/                     # Config has DB sqlexec.DB; no DSN field
    └── nil-db-cli-usage/                    # cli.Run status with nil DB → exit 2
```

**Significance order:** surface (interface | wrap | exec | query | query-row |
close | config) → outcome class (happy vs error/empty) → concrete SQL/table.

Offline leaves (no MySQL): `interface/methods-present`, `config/*`.
Live-MySQL leaves skip when DSN unreachable.

## Test Index

| Leaf | Description |
|------|-------------|
| `interface/methods-present` | Package exports `DB`, `Result`, `Rows`, `Row`, `Wrap`; method sets match locked API |
| `wrap/returns-db` | `Wrap(sqlDB)` returns non-nil `DB`; `Exec`/`QueryRow` of `SELECT 1` succeeds |
| `exec/success` | `Exec` CREATE TABLE + INSERT → `RowsAffected() >= 1`; cleanup DROP |
| `exec/bad-sql` | `Exec` of invalid SQL returns error |
| `query/multi-row` | Seed 2 rows; `Query` scans both in order |
| `query/empty` | `Query` matching nothing → no `Next`, `Err` nil, `Close` ok |
| `query-row/one-row` | `QueryRow` + `Scan` returns the seeded value |
| `query-row/no-rows` | `QueryRow` + `Scan` → `errors.Is(err, sql.ErrNoRows)` |
| `close/after-close-exec-errors` | After `Close()`, `Exec` returns error |
| `config/db-field-no-dsn` | `migrate.Config` has `DB` of type `sqlexec.DB`; **no** `DSN` field |
| `config/nil-db-cli-usage` | `cli.Run(cfg{DB:nil}, ["status"])` → exit **2**; Error about DB/missing |

## How to Run

```sh
cd /Users/xhd2015/Projects/xhd2015/mysql-migrate
# optional: export MIGRATE_MYSQL_DSN='user:pass@tcp(host:port)/db?charset=utf8mb4&parseTime=True'
doctest vet ./tests/sqlexec
doctest test ./tests/sqlexec
```

Live leaves need MySQL at the resolved DSN (default `localhost:9306` /
`lifespan_db`). Offline leaves always run. Live leaves **skip** when the DSN is
not reachable so pure-unit trees stay usable offline.

Classic TDD: importing `migrate/sqlexec` fails compile until implementer lands
the package. Leaves must fail (compile or assertion RED) until:

```text
migrate/sqlexec
```

### Locked API (implementer)

```go
package sqlexec

import (
    "context"
    "database/sql"
)

// DB is the context-first SQL handle used by logrepo/cli (never *sql.DB).
type DB interface {
    Exec(ctx context.Context, query string, args ...any) (Result, error)
    Query(ctx context.Context, query string, args ...any) (Rows, error)
    QueryRow(ctx context.Context, query string, args ...any) Row
    Close() error
}

// Result is the outcome of Exec (mirrors database/sql.Result surface used by migrate).
type Result interface {
    LastInsertId() (int64, error)
    RowsAffected() (int64, error)
}

// Rows is a multi-row cursor from Query.
type Rows interface {
    Next() bool
    Scan(dest ...any) error
    Close() error
    Err() error
}

// Row is a single-row accessor from QueryRow.
type Row interface {
    Scan(dest ...any) error
}

// Wrap adapts *sql.DB into DB by forwarding to *Context methods.
// Does not open connections; does not take a DSN.
func Wrap(db *sql.DB) DB
```

### Config.DB-only (implementer)

```go
package migrate

import "github.com/xhd2015/mysql-migrate/migrate/sqlexec"

type Config struct {
    DB            sqlexec.DB // required for DB subcommands; nil → usage error in cli
    MigrationsDir string
    ProgramName   string
    AppliedBy     string
}
// No DSN field — library never sql.Open.
```

### Harness notes for existing trees

| Tree | P1 harness expectation |
|------|------------------------|
| `tests/scaffold` | `Config` asserts `DB` field + identity strings; **no** `DSN` field |
| `tests/cli` | `buildConfig` opens MySQL only in harness, `cfg.DB = sqlexec.Wrap(sqlDB)`; leaves never put DSN on Config; `usage/status-missing-db` for nil DB |
| `tests/logrepo` | `Run` uses `sqlexec.Wrap` and passes `sqlexec.DB` into logrepo APIs |
| `cmd` | Out of scope for P1 (P2 may wire DSN → Wrap → Config.DB) |

```go
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/xhd2015/mysql-migrate/cli"
	"github.com/xhd2015/mysql-migrate/migrate"
	"github.com/xhd2015/mysql-migrate/migrate/sqlexec"
)

// defaultLocalDSN is the lifelog/local-dev MySQL DSN used when MIGRATE_MYSQL_DSN
// is unset. multiStatements not required for sqlexec leaves.
const defaultLocalDSN = "lf:Xpassword@tcp(localhost:9306)/lifespan_db?charset=utf8mb4&parseTime=True"

// Request drives one sqlexec / Config.DB scenario.
type Request struct {
	// Op is the scenario dispatch key:
	// interface | wrap | exec_ok | exec_err | query_multi | query_empty |
	// query_row_ok | query_row_none | close | config_fields | cli_nil_db
	Op string

	// Table is an isolated MySQL table name for live leaves (session-prefixed).
	Table string

	// Query is optional SQL override (exec_err may set a fixed bad statement).
	Query string

	// SeedValues are integer payloads for multi-row / one-row seeds.
	SeedValues []int64

	// MigrationsDir for cli_nil_db (must be non-empty so failure is about DB).
	MigrationsDir string
}

// Response holds observed sqlexec / Config / cli outcomes.
type Response struct {
	// Interface / config_fields (compile-time shape probes).
	HasDBType       bool
	HasResultType   bool
	HasRowsType     bool
	HasRowType      bool
	HasWrapFunc     bool
	DBMethodsOK     bool // Exec, Query, QueryRow, Close present on interface
	ResultMethodsOK bool
	RowsMethodsOK   bool
	RowMethodsOK    bool
	ConfigHasDB     bool
	ConfigDBIsIface bool // Config.DB type is interface assignable from sqlexec.DB
	ConfigHasDSN    bool // must be false after P1
	ConfigDBNilZero bool // zero Config.DB is nil

	// Wrap / live ops.
	WrapNonNil bool
	ExecOK     bool
	RowsAffected int64
	QueryCount int
	Scanned    []int64
	QueryEmpty bool // Next never true and Err nil
	ScanValue  int64
	CloseErr   error // error from Close() itself (expect nil)
	PostCloseExecErr bool

	// Error paths.
	OpErr      error  // primary op error (exec_err, query_row_none, …)
	OpErrIsNoRows bool // errors.Is(OpErr, sql.ErrNoRows)

	// cli_nil_db
	ExitCode int
	Stdout   string
	Stderr   string
}

// resolveDSN returns MIGRATE_MYSQL_DSN if set, else defaultLocalDSN.
func resolveDSN() string {
	if dsn := strings.TrimSpace(os.Getenv("MIGRATE_MYSQL_DSN")); dsn != "" {
		return dsn
	}
	return defaultLocalDSN
}

// openSQL opens and pings MySQL. Caller closes.
func openSQL(t *testing.T) (*sql.DB, error) {
	t.Helper()
	db, err := sql.Open("mysql", resolveDSN())
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping MySQL: %w", err)
	}
	return db, nil
}

// wrapOrFail opens SQL, wraps via sqlexec.Wrap, returns (db, raw, cleanup, err).
// cleanup closes the sqlexec.DB (preferred) or raw on error paths.
func wrapOrFail(t *testing.T) (sqlexec.DB, *sql.DB, func(), error) {
	t.Helper()
	raw, err := openSQL(t)
	if err != nil {
		return nil, nil, func() {}, err
	}
	db := sqlexec.Wrap(raw)
	cleanup := func() {
		if db != nil {
			_ = db.Close()
			return
		}
		_ = raw.Close()
	}
	return db, raw, cleanup, nil
}

func ctxTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// probeInterface records whether the locked types/methods exist via compile-time
// usage and light reflection on method names of the interface types.
func probeInterface(resp *Response) {
	// Compile-time: names must exist.
	var _ sqlexec.DB
	var _ sqlexec.Result
	var _ sqlexec.Rows
	var _ sqlexec.Row
	_ = sqlexec.Wrap

	resp.HasDBType = true
	resp.HasResultType = true
	resp.HasRowsType = true
	resp.HasRowType = true
	resp.HasWrapFunc = true

	// Method sets via reflect on interface types.
	dbT := reflect.TypeOf((*sqlexec.DB)(nil)).Elem()
	resp.DBMethodsOK = hasMethods(dbT, "Exec", "Query", "QueryRow", "Close")

	resT := reflect.TypeOf((*sqlexec.Result)(nil)).Elem()
	resp.ResultMethodsOK = hasMethods(resT, "LastInsertId", "RowsAffected")

	rowsT := reflect.TypeOf((*sqlexec.Rows)(nil)).Elem()
	resp.RowsMethodsOK = hasMethods(rowsT, "Next", "Scan", "Close", "Err")

	rowT := reflect.TypeOf((*sqlexec.Row)(nil)).Elem()
	resp.RowMethodsOK = hasMethods(rowT, "Scan")
}

func hasMethods(iface reflect.Type, names ...string) bool {
	if iface == nil || iface.Kind() != reflect.Interface {
		return false
	}
	for _, name := range names {
		if _, ok := iface.MethodByName(name); !ok {
			return false
		}
	}
	return true
}

func probeConfigFields(resp *Response) {
	typ := reflect.TypeOf(migrate.Config{})
	if f, ok := typ.FieldByName("DB"); ok {
		resp.ConfigHasDB = true
		// Field type should be the sqlexec.DB interface.
		dbIface := reflect.TypeOf((*sqlexec.DB)(nil)).Elem()
		resp.ConfigDBIsIface = f.Type == dbIface
	}
	_, resp.ConfigHasDSN = typ.FieldByName("DSN")

	var zero migrate.Config
	resp.ConfigDBNilZero = zero.DB == nil
}

// Run exercises one sqlexec / Config.DB scenario selected by req.Op.
// Classic RED until migrate/sqlexec exists and Config loses DSN.
func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	if req.Op == "" {
		return nil, fmt.Errorf("empty Op")
	}

	resp := &Response{}

	switch req.Op {
	case "interface":
		probeInterface(resp)
		return resp, nil

	case "config_fields":
		probeInterface(resp) // needs package present
		probeConfigFields(resp)
		return resp, nil

	case "cli_nil_db":
		// Offline usage path: nil DB must fail before any open.
		cfg := migrate.Config{
			DB:            nil,
			MigrationsDir: req.MigrationsDir,
			ProgramName:   "mysql-migrate",
			AppliedBy:     "sqlexec-doctest",
		}
		// Capture is minimal — reuse os redirection only if needed.
		// Prefer calling cli.Run directly; tests capture via helpers in SETUP when present.
		// Here we call Run and only observe exit code; stdout/stderr may be empty if
		// implementation writes elsewhere — Assert checks exit + optional tokens via
		// re-invocation with pipes in this case.
		stdout, stderr, code := captureCLI(t, cfg, []string{"status"})
		resp.ExitCode = code
		resp.Stdout = stdout
		resp.Stderr = stderr
		return resp, nil

	case "wrap", "exec_ok", "exec_err", "query_multi", "query_empty",
		"query_row_ok", "query_row_none", "close":
		db, raw, cleanup, err := wrapOrFail(t)
		if err != nil {
			return nil, err
		}
		defer cleanup()
		resp.WrapNonNil = db != nil
		if db == nil {
			return resp, fmt.Errorf("sqlexec.Wrap returned nil")
		}

		switch req.Op {
		case "wrap":
			ctx, cancel := ctxTimeout()
			defer cancel()
			// Prove the wrapped handle can talk to MySQL.
			row := db.QueryRow(ctx, "SELECT 1")
			var n int
			if err := row.Scan(&n); err != nil {
				resp.OpErr = err
				return resp, nil
			}
			resp.ScanValue = int64(n)
			resp.ExecOK = n == 1
			return resp, nil

		case "exec_ok":
			return runExecOK(t, db, raw, req, resp)

		case "exec_err":
			ctx, cancel := ctxTimeout()
			defer cancel()
			q := req.Query
			if q == "" {
				q = "THIS IS NOT VALID SQL !!!"
			}
			_, err := db.Exec(ctx, q)
			resp.OpErr = err
			return resp, nil

		case "query_multi":
			return runQueryMulti(t, db, raw, req, resp)

		case "query_empty":
			return runQueryEmpty(t, db, raw, req, resp)

		case "query_row_ok":
			return runQueryRowOK(t, db, raw, req, resp)

		case "query_row_none":
			return runQueryRowNone(t, db, raw, req, resp)

		case "close":
			// Close the facade; further Exec must error.
			if err := db.Close(); err != nil {
				resp.CloseErr = err
			}
			// Prevent double-close of raw in cleanup: nil out by closing only once.
			// cleanup still calls db.Close — second Close on sql.DB is safe (returns err).
			ctx, cancel := ctxTimeout()
			defer cancel()
			_, err := db.Exec(ctx, "SELECT 1")
			resp.PostCloseExecErr = err != nil
			resp.OpErr = err
			// Mark db closed: cleanup will Close again (idempotent-ish).
			return resp, nil
		}

	default:
		return nil, fmt.Errorf("unknown op %q", req.Op)
	}
	return resp, nil
}

func runExecOK(t *testing.T, db sqlexec.DB, raw *sql.DB, req *Request, resp *Response) (*Response, error) {
	t.Helper()
	table := req.Table
	if table == "" {
		return nil, fmt.Errorf("exec_ok requires Table")
	}
	ctx, cancel := ctxTimeout()
	defer cancel()

	// Best-effort cleanup via raw in case Wrap Close is wrong during RED/GREEN.
	defer func() {
		_, _ = raw.ExecContext(context.Background(), "DROP TABLE IF EXISTS "+table)
	}()

	if _, err := db.Exec(ctx, "DROP TABLE IF EXISTS "+table); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	if _, err := db.Exec(ctx, "CREATE TABLE "+table+" (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, v BIGINT NOT NULL) ENGINE=InnoDB"); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	res, err := db.Exec(ctx, "INSERT INTO "+table+" (v) VALUES (?)", int64(42))
	if err != nil {
		resp.OpErr = err
		return resp, nil
	}
	n, err := res.RowsAffected()
	if err != nil {
		resp.OpErr = err
		return resp, nil
	}
	resp.RowsAffected = n
	resp.ExecOK = n >= 1
	return resp, nil
}

func seedTable(t *testing.T, db sqlexec.DB, raw *sql.DB, table string, values []int64) error {
	t.Helper()
	ctx, cancel := ctxTimeout()
	defer cancel()
	if _, err := db.Exec(ctx, "DROP TABLE IF EXISTS "+table); err != nil {
		return err
	}
	if _, err := db.Exec(ctx, "CREATE TABLE "+table+" (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, v BIGINT NOT NULL) ENGINE=InnoDB"); err != nil {
		return err
	}
	for _, v := range values {
		if _, err := db.Exec(ctx, "INSERT INTO "+table+" (v) VALUES (?)", v); err != nil {
			return err
		}
	}
	// Register cleanup on the test.
	t.Cleanup(func() {
		_, _ = raw.ExecContext(context.Background(), "DROP TABLE IF EXISTS "+table)
	})
	return nil
}

func runQueryMulti(t *testing.T, db sqlexec.DB, raw *sql.DB, req *Request, resp *Response) (*Response, error) {
	t.Helper()
	if req.Table == "" {
		return nil, fmt.Errorf("query_multi requires Table")
	}
	vals := req.SeedValues
	if len(vals) == 0 {
		vals = []int64{10, 20}
	}
	if err := seedTable(t, db, raw, req.Table, vals); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	rows, err := db.Query(ctx, "SELECT v FROM "+req.Table+" ORDER BY id")
	if err != nil {
		resp.OpErr = err
		return resp, nil
	}
	defer rows.Close()
	var got []int64
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			resp.OpErr = err
			return resp, nil
		}
		got = append(got, v)
	}
	if err := rows.Err(); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	resp.Scanned = got
	resp.QueryCount = len(got)
	return resp, nil
}

func runQueryEmpty(t *testing.T, db sqlexec.DB, raw *sql.DB, req *Request, resp *Response) (*Response, error) {
	t.Helper()
	if req.Table == "" {
		return nil, fmt.Errorf("query_empty requires Table")
	}
	if err := seedTable(t, db, raw, req.Table, nil); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	rows, err := db.Query(ctx, "SELECT v FROM "+req.Table+" WHERE v = ?", int64(-999999))
	if err != nil {
		resp.OpErr = err
		return resp, nil
	}
	defer rows.Close()
	if rows.Next() {
		resp.QueryEmpty = false
		resp.QueryCount = 1
		return resp, nil
	}
	if err := rows.Err(); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	resp.QueryEmpty = true
	resp.QueryCount = 0
	return resp, nil
}

func runQueryRowOK(t *testing.T, db sqlexec.DB, raw *sql.DB, req *Request, resp *Response) (*Response, error) {
	t.Helper()
	if req.Table == "" {
		return nil, fmt.Errorf("query_row_ok requires Table")
	}
	vals := req.SeedValues
	if len(vals) == 0 {
		vals = []int64{77}
	}
	if err := seedTable(t, db, raw, req.Table, vals); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	row := db.QueryRow(ctx, "SELECT v FROM "+req.Table+" ORDER BY id LIMIT 1")
	var v int64
	if err := row.Scan(&v); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	resp.ScanValue = v
	return resp, nil
}

func runQueryRowNone(t *testing.T, db sqlexec.DB, raw *sql.DB, req *Request, resp *Response) (*Response, error) {
	t.Helper()
	if req.Table == "" {
		return nil, fmt.Errorf("query_row_none requires Table")
	}
	if err := seedTable(t, db, raw, req.Table, nil); err != nil {
		resp.OpErr = err
		return resp, nil
	}
	ctx, cancel := ctxTimeout()
	defer cancel()
	row := db.QueryRow(ctx, "SELECT v FROM "+req.Table+" WHERE v = ?", int64(-999999))
	var v int64
	err := row.Scan(&v)
	resp.OpErr = err
	resp.OpErrIsNoRows = errors.Is(err, sql.ErrNoRows)
	return resp, nil
}

// captureCLI runs cli.Run with stdio redirected to buffers.
func captureCLI(t *testing.T, cfg migrate.Config, args []string) (stdout, stderr string, code int) {
	t.Helper()
	// Use os pipes for fidelity with other trees.
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	rErr, wErr, err := os.Pipe()
	if err != nil {
		_ = rOut.Close()
		_ = wOut.Close()
		t.Fatalf("pipe stderr: %v", err)
	}
	oldOut, oldErr, oldIn := os.Stdout, os.Stderr, os.Stdin
	os.Stdout, os.Stderr = wOut, wErr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		os.Stdout, os.Stderr, os.Stdin = oldOut, oldErr, oldIn
		t.Fatalf("open /dev/null: %v", err)
	}
	os.Stdin = devNull

	outCh := make(chan string, 1)
	errCh := make(chan string, 1)
	go func() {
		var b strings.Builder
		buf := make([]byte, 4096)
		for {
			n, e := rOut.Read(buf)
			if n > 0 {
				b.Write(buf[:n])
			}
			if e != nil {
				break
			}
		}
		outCh <- b.String()
	}()
	go func() {
		var b strings.Builder
		buf := make([]byte, 4096)
		for {
			n, e := rErr.Read(buf)
			if n > 0 {
				b.Write(buf[:n])
			}
			if e != nil {
				break
			}
		}
		errCh <- b.String()
	}()

	code = cli.Run(cfg, args)

	_ = wOut.Close()
	_ = wErr.Close()
	stdout = <-outCh
	stderr = <-errCh
	_ = rOut.Close()
	_ = rErr.Close()
	_ = devNull.Close()
	os.Stdout, os.Stderr, os.Stdin = oldOut, oldErr, oldIn
	return stdout, stderr, code
}
```
