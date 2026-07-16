# mysql-migrate — README / docs polish (P7)

Standalone suite polish: the module ships a **root `README.md`** that explains
purpose, install, CLI usage (`--dsn` / `--dir`), subcommands, env vars, and how
to run the doctest suite. This tree is **docs-only** — it reads
`README.md` at the module root and asserts required phrases.

Classic RED until `README.md` exists and contains the locked phrases.
No production package code is required for this tree (no `cli` / `migrate`
imports). No lifelog package imports.

Standalone doctest root under `tests/docs/` so scaffold / inventory / plan /
logrepo / cli / cmd trees stay independent.

# DSN (Domain Specific Notion)

The **module README** is the **human front door** for the **mysql-migrate**
repository. An **operator** or **contributor** opens **`README.md`** at the
module root (next to `go.mod`). The README **documents**:

1. **Purpose** — what the tool does (MySQL migration operator).
2. **Install** — how to obtain the binary or module
   (`go install` / module path `github.com/xhd2015/mysql-migrate`).
3. **CLI usage** — global flags **`--dsn`** and **`--dir`** and example
   invocations of the **`mysql-migrate`** binary.
4. **Subcommands** — `status`, `plan`, `apply`, `mark-done`, `mark-failed`,
   `note`, `allow-retry`.
5. **Environment** — optional fallbacks **`MIGRATE_MYSQL_DSN`** and
   **`MIGRATE_MYSQL_DIR`** when flags are omitted.
6. **Running tests** — how to run **doctests** (`doctest vet` / `doctest test`)
   and/or `go test` against this repo.

This test tree is a **reader**: it **loads** `README.md` from the module root
and **checks** that required phrases for each documentation topic are present.
Missing README or missing phrases fail the leaf (implementer writes the
README to turn RED → GREEN). The tree does **not** execute the CLI or open
MySQL.

## Version

0.0.2

## Decision Tree

Split on **README section topic** (documentation contract surface):

```
tests/docs/                              [Request{RequiredPhrases}]
│                                        Run: read module-root README.md
└── readme-section/                      # all leaves assert README body text
    ├── purpose/                         # tool name + MySQL migration purpose
    ├── install/                         # go install / module path
    ├── cli-usage/                       # --dsn, --dir, mysql-migrate usage
    ├── subcommands/                     # status plan apply mark-* note allow-retry
    ├── env-vars/                        # MIGRATE_MYSQL_DSN, MIGRATE_MYSQL_DIR
    └── run-doctests/                    # doctest test / how to run suite
```

Siblings are MECE over P7 README exit criteria (purpose, install, CLI flags,
subcommands, env, doctest runbook). No CLI/DB behavior here (covered by
`tests/cli` and `tests/cmd`).

**Significance order:** document = README (only surface) → section topic
(purpose / install / flags / commands / env / test runbook).

## Test Case Index

| # | Path | Preconditions | Expected |
|---|------|---------------|----------|
| 1 | `readme-section/purpose/` | README exists at module root | contains `mysql-migrate` and migration purpose phrasing |
| 2 | `readme-section/install/` | README exists | documents install via `go install` and module path |
| 3 | `readme-section/cli-usage/` | README exists | documents `--dsn` and `--dir` |
| 4 | `readme-section/subcommands/` | README exists | lists all seven subcommands |
| 5 | `readme-section/env-vars/` | README exists | mentions `MIGRATE_MYSQL_DSN` and `MIGRATE_MYSQL_DIR` |
| 6 | `readme-section/run-doctests/` | README exists | documents running `doctest test` (and suite location) |

## How to Run

```sh
cd /Users/xhd2015/Projects/xhd2015/mysql-migrate
doctest vet ./tests/docs
doctest test ./tests/docs
# implementer exit check after README is written:
# go test ./...
# doctest test ./...
```

Classic RED: without `README.md` (or without the locked phrases), every leaf
fails Assert with a clear missing-file or missing-phrase message.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Request names the phrases this leaf requires in README.md.
// Leaves set RequiredPhrases in Setup; root Run only reads the file.
type Request struct {
	// RequiredPhrases are case-sensitive substrings that must appear in README.md.
	RequiredPhrases []string

	// Label is a short section name used in failure messages (e.g. "purpose").
	Label string
}

// Response is the observed module-root README content (or absence).
type Response struct {
	ModuleRoot string
	READMEPath string
	Exists     bool
	Content    string
	// Size is len(Content) when the file was read successfully.
	Size int
}

// Run loads README.md from the module root (two levels above this tree).
// Missing file is an observed outcome (Exists=false), not a harness error.
func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	moduleRoot := filepath.Clean(filepath.Join(DOCTEST_ROOT, "..", ".."))
	readmePath := filepath.Join(moduleRoot, "README.md")

	resp := &Response{
		ModuleRoot: moduleRoot,
		READMEPath: readmePath,
	}

	data, err := os.ReadFile(readmePath)
	if err != nil {
		if os.IsNotExist(err) {
			resp.Exists = false
			return resp, nil
		}
		return nil, fmt.Errorf("read %s: %w", readmePath, err)
	}

	resp.Exists = true
	resp.Content = string(data)
	resp.Size = len(data)
	return resp, nil
}
```
