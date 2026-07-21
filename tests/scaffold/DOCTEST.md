# mysql-migrate — module scaffold + Config (P1, DB-only)

Standalone module scaffold: public `Config` type under package `migrate`,
module path `github.com/xhd2015/mysql-migrate`, and a buildable package layout.
Stubs for inventory/plan/logrepo/cli/cmd are allowed; no migration logic yet.

P1 sqlexec contract: **Config carries `DB sqlexec.DB`**, not a DSN string.
Opening connections is outside the library (callers use `sqlexec.Wrap`).

# DSN (Domain Specific Notion)

**Participants**

- **Module** — Go module `github.com/xhd2015/mysql-migrate` rooted at the
  repository (has `go.mod`).
- **migrate package** — import path
  `github.com/xhd2015/mysql-migrate/migrate`; holds the public `Config` type
  (file such as `migrate/config.go`).
- **Config** — connection and identity settings for later migration commands:
  `DB` (`sqlexec.DB`), `MigrationsDir`, `ProgramName`, `AppliedBy`. **No DSN
  field**.
- **sqlexec package** — `github.com/xhd2015/mysql-migrate/migrate/sqlexec`
  provides the `DB` interface type referenced by Config (may be stubbed until
  full implementer pass; field type must still compile).
- **Stub packages** (optional at scaffold) — `migrate/inventory`, `migrate/plan`,
  `migrate/logrepo`, `cli`, `cmd/mysql-migrate`; may be empty packages so
  `go build ./...` still succeeds.
- **Caller** — any consumer (tests, future CLI) that imports `migrate` and
  constructs `Config`.

**Behaviors**

- Caller imports `github.com/xhd2015/mysql-migrate/migrate` and obtains type
  `Config` without build errors.
- Zero-value `Config{}` is valid; `DB` is nil; string fields default to empty.
- Populated `Config` stores `DB` and the three string fields as set.
- From the module root, `go build ./...` exits 0 (stubs allowed).

## Version

0.0.2

## Decision Tree

Split on **verification subject** (what exit-criterion surface we check):

```
tests/scaffold/
├── config-fields-exist/   # Config field surface: zero + populated (DB, not DSN)
├── module-importable/     # import path + package identity via reflect
└── packages-build/        # go build ./... at module root
```

Siblings are MECE over scaffold exit criteria: type shape, importability, and
module buildability. No inventory/plan/apply behavior here (later phases).

## Test Case Index

| # | Path | Preconditions | Expected |
|---|------|---------------|----------|
| 1 | `config-fields-exist/` | construct zero-value and populated `migrate.Config` | `DB` nil/zero; strings readable; **no DSN field** |
| 2 | `module-importable/` | import `migrate` package | `PkgPath` is `github.com/xhd2015/mysql-migrate/migrate`; type name `Config` |
| 3 | `packages-build/` | module root = repo (two levels above this tree) | `go build ./...` exit 0 |

## How to Run

```sh
cd /Users/xhd2015/Projects/xhd2015/mysql-migrate
doctest vet ./tests/scaffold
doctest test ./tests/scaffold
# implementer exit check:
go build ./...
```

Classic RED: until `migrate.Config` has `DB sqlexec.DB` (and no DSN), the
generated tests fail to compile or assert RED. Full `sqlexec` behavior is
covered under `tests/sqlexec/`.

```go
import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate"
	"github.com/xhd2015/mysql-migrate/migrate/sqlexec"
	"github.com/xhd2015/doctest/session"
)

// Mode selects which scaffold surface Run exercises.
// "config-fields-exist" | "module-importable" | "packages-build"
type Request struct {
	Mode string

	// config-fields-exist sample values (DB stays nil in scaffold; type only)
	SampleMigrationsDir string
	SampleProgramName   string
	SampleAppliedBy     string
}

type Response struct {
	// config-fields-exist
	ZeroDBIsNil         bool
	ZeroMigrationsDir   string
	ZeroProgramName     string
	ZeroAppliedBy       string
	PopulatedDBIsNil    bool // scaffold leaves DB nil; still valid Config
	MigrationsDir       string
	ProgramName         string
	AppliedBy           string
	HasDBField          bool
	DBFieldIsSqlexecDB  bool
	HasDSNField         bool // must be false

	// module-importable
	PkgPath  string
	TypeName string

	// packages-build
	BuildExitCode int
	BuildStdout   string
	BuildStderr   string
	ModuleRoot    string
}

// Run exercises the scaffold surface selected by req.Mode.
// Classic RED until Config.DB exists and DSN is removed.
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	switch req.Mode {
	case "config-fields-exist":
		var zero migrate.Config
		populated := migrate.Config{
			DB:            nil, // type presence; live Wrap is tests/sqlexec
			MigrationsDir: req.SampleMigrationsDir,
			ProgramName:   req.SampleProgramName,
			AppliedBy:     req.SampleAppliedBy,
		}
		typ := reflect.TypeOf(migrate.Config{})
		dbField, hasDB := typ.FieldByName("DB")
		_, hasDSN := typ.FieldByName("DSN")
		dbIface := reflect.TypeOf((*sqlexec.DB)(nil)).Elem()
		return &Response{
			ZeroDBIsNil:        zero.DB == nil,
			ZeroMigrationsDir:  zero.MigrationsDir,
			ZeroProgramName:    zero.ProgramName,
			ZeroAppliedBy:      zero.AppliedBy,
			PopulatedDBIsNil:   populated.DB == nil,
			MigrationsDir:      populated.MigrationsDir,
			ProgramName:        populated.ProgramName,
			AppliedBy:          populated.AppliedBy,
			HasDBField:         hasDB,
			DBFieldIsSqlexecDB: hasDB && dbField.Type == dbIface,
			HasDSNField:        hasDSN,
		}, nil

	case "module-importable":
		typ := reflect.TypeOf(migrate.Config{})
		return &Response{
			PkgPath:  typ.PkgPath(),
			TypeName: typ.Name(),
		}, nil

	case "packages-build":
		moduleRoot := filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
		cmd := exec.Command("go", "build", "./...")
		cmd.Dir = moduleRoot
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		runErr := cmd.Run()
		resp := &Response{
			BuildStdout: stdout.String(),
			BuildStderr: stderr.String(),
			ModuleRoot:  moduleRoot,
		}
		if runErr != nil {
			if exitErr, ok := runErr.(*exec.ExitError); ok {
				resp.BuildExitCode = exitErr.ExitCode()
				// Build failure is an observed outcome for Assert, not a harness error.
				return resp, nil
			}
			return nil, fmt.Errorf("go build ./... failed: %w\nstdout:\n%s\nstderr:\n%s",
				runErr, stdout.String(), stderr.String())
		}
		resp.BuildExitCode = 0
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown mode %q", req.Mode)
	}
}
```
