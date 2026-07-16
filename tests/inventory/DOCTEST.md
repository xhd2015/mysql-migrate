# mysql-migrate — file inventory (P2)

Pure library tests for discovering, parsing, sorting, and hashing append-only
MySQL migration SQL files. No database, no CLI apply, no log table.

Target package (Classic RED until implementer ports logic):

```text
github.com/xhd2015/mysql-migrate/migrate/inventory
```

# DSN (Domain Specific Notion)

The **migration inventory** is a pure on-disk reader of the **migrations
directory**. An **author** drops **migration files** at the top level using the
filename grammar `YYYY-MM-DD-NN[-[EXACTLY-ONCE]]-<slug>.sql`. The inventory
**parses** basenames into structured metadata (`migration_id`, date, seq, slug,
`ExactlyOnce`), **lists** only top-level grammar-matching `.sql` files
(lexicographic by full filename), **ignores** unrelated top-level junk and
everything under subdirectories, and **hashes** raw file bytes with SHA-256
(lowercase hex). Names that loosely look like migrations (`YYYY-MM-DD-*.sql`)
but fail the grammar surface as **errors** on list; plain unrelated files do
not. Hash is content-only: same bytes → same digest.

**Participants**

- **Author** — places migration SQL files under a migrations directory.
- **Inventory package** — `ParseFileName`, `ListDir`, `HashFile` (no MySQL).
- **Migration file** — top-level `.sql` matching the grammar; metadata includes
  `ID`, `FileName`, `Path`, `Date`, `Seq`, `Slug`, `ExactlyOnce`, `ContentSHA256`.
- **Migrations directory** — top-level scan only; subdirectories are never listed.

**Behaviors**

- `ParseFileName(name)` → structured `MigrationFile` or error (no Path/hash).
- `ListDir(dir)` → sorted top-level matches with Path + ContentSHA256 filled.
- `HashFile(path)` → 64-char lowercase hex SHA-256 of raw bytes.
- Grammar lookalikes on list → error; unrelated junk → ignored.

## Version

0.0.2

## Decision Tree

Split on **operation** (`Op`: parse | list | hash) — largest behavioral fork:

```
tests/inventory/                          [Request{Op, FileName, Dir, Path, ...}]
│                                         Run: inventory.ParseFileName | ListDir | HashFile
├── parse/                                # ParseFileName basename → metadata
│   ├── valid/
│   │   ├── simple/                       # date-NN-slug, ExactlyOnce=false
│   │   ├── exactly-once/                 # middle [EXACTLY-ONCE] token
│   │   └── stem-without-extension/       # basename without .sql still parses
│   └── invalid/
│       └── rejects-bad-names/            # table of invalid basenames → error
├── list/                                 # ListDir top-level scan + hash fill
│   ├── sort-by-filename/                 # out-of-order files → name ASC
│   ├── ignores-non-matching-and-subdirs/ # junk + nested .sql ignored
│   ├── includes-sha256/                  # ContentSHA256 == HashFile(path)
│   └── rejects-grammar-lookalikes/       # YYYY-MM-DD-*.sql bad grammar → error
└── hash/                                 # HashFile content digest
    ├── stable-same-bytes/                # same path/content → same hex
    └── differs-on-content-change/        # different bytes → different hex
```

**Significance order:** `Op` (parse | list | hash) → validity / list outcome /
content relation → concrete filename or fixture details.

Siblings at each level are MECE for the split factor (parse vs list vs hash;
valid vs invalid parse; sort / ignore / hash-fill / lookalike-error for list;
stable vs differs for hash).

## Test Case Index

| # | Path | Description |
|---|------|-------------|
| 1 | `parse/valid/simple` | `2026-07-16-01-create-t-channel.sql` → id, date, seq=1, ExactlyOnce=false, slug |
| 2 | `parse/valid/exactly-once` | Middle `[EXACTLY-ONCE]` → ExactlyOnce=true, slug after marker |
| 3 | `parse/valid/stem-without-extension` | Stem without `.sql` parses same metadata |
| 4 | `parse/invalid/rejects-bad-names` | Several invalid basenames each return error |
| 5 | `list/sort-by-filename` | Fixture dir out of order → sorted by FileName ASC |
| 6 | `list/ignores-non-matching-and-subdirs` | README, nested sql ignored; only valid top-level |
| 7 | `list/includes-sha256` | Listed `ContentSHA256` matches `HashFile` |
| 8 | `list/rejects-grammar-lookalikes` | Top-level `YYYY-MM-DD-*.sql` failing grammar → list error |
| 9 | `hash/stable-same-bytes` | Two HashFile calls on same bytes → identical 64-char hex |
| 10 | `hash/differs-on-content-change` | Different content → different digests |

## How to Run

```sh
cd /Users/xhd2015/Projects/xhd2015/mysql-migrate
doctest vet ./tests/inventory
doctest test ./tests/inventory
```

Classic RED: package is stub-only (`migrate/inventory/doc.go`) until implementer
ports `ParseFileName`, `ListDir`, `HashFile`, and `MigrationFile`. Do not implement
here — designer owns tests only.

```go
import (
	"fmt"
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
)

// Request drives one inventory operation.
type Request struct {
	// Op is the dispatch key: parse | parse_expect_errors | list | hash.
	Op string

	// ParseFileName input (with or without .sql).
	FileName string

	// InvalidNames is a batch of basenames that must each fail ParseFileName
	// (Op=parse_expect_errors).
	InvalidNames []string

	// Dir is the migrations directory for ListDir.
	Dir string

	// Path / PathB are file paths for HashFile (PathB optional).
	Path  string
	PathB string

	// FixtureFiles maps top-level relative name → content written under Dir
	// before list/hash ops (created in leaf Setup via helpers).
	FixtureFiles map[string]string

	// FixtureNested maps relative path under Dir (may include subdirs) → content.
	FixtureNested map[string]string
}

// FileView is a plain snapshot of inventory.MigrationFile for assertions.
type FileView struct {
	ID            string
	FileName      string
	Path          string
	ExactlyOnce   bool
	Date          string
	Seq           int
	Slug          string
	ContentSHA256 string
}

// Response holds outcomes of parse / list / hash.
type Response struct {
	File  *FileView
	Files []FileView

	Hash  string
	HashB string

	// ParseErrors[i] is the error text for InvalidNames[i]; empty string means
	// ParseFileName unexpectedly succeeded for that name.
	ParseErrors []string
}

// Run exercises the inventory surface selected by req.Op.
// Classic RED until migrate/inventory exports MigrationFile + ParseFileName + ListDir + HashFile.
func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}

	// Materialise fixtures into Dir when present (list/hash leaves).
	if req.Dir != "" && (len(req.FixtureFiles) > 0 || len(req.FixtureNested) > 0) {
		if err := materializeFixtures(t, req.Dir, req.FixtureFiles, req.FixtureNested); err != nil {
			return nil, err
		}
	}

	switch req.Op {
	case "parse":
		mf, err := inventory.ParseFileName(req.FileName)
		if err != nil {
			return &Response{}, err
		}
		v := toFileView(mf)
		return &Response{File: &v}, nil

	case "parse_expect_errors":
		resp := &Response{ParseErrors: make([]string, 0, len(req.InvalidNames))}
		for _, name := range req.InvalidNames {
			_, err := inventory.ParseFileName(name)
			if err != nil {
				resp.ParseErrors = append(resp.ParseErrors, err.Error())
			} else {
				resp.ParseErrors = append(resp.ParseErrors, "")
			}
		}
		return resp, nil

	case "list":
		files, err := inventory.ListDir(req.Dir)
		if err != nil {
			return &Response{}, err
		}
		out := make([]FileView, 0, len(files))
		for _, f := range files {
			out = append(out, toFileView(f))
		}
		return &Response{Files: out}, nil

	case "hash":
		h, err := inventory.HashFile(req.Path)
		if err != nil {
			return &Response{}, err
		}
		resp := &Response{Hash: h}
		if req.PathB != "" {
			h2, err := inventory.HashFile(req.PathB)
			if err != nil {
				return &Response{Hash: h}, err
			}
			resp.HashB = h2
		} else {
			// Second call on same path — stability check helper.
			h2, err := inventory.HashFile(req.Path)
			if err != nil {
				return &Response{Hash: h}, err
			}
			resp.HashB = h2
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown op %q", req.Op)
	}
}

func toFileView(mf inventory.MigrationFile) FileView {
	return FileView{
		ID:            mf.ID,
		FileName:      mf.FileName,
		Path:          mf.Path,
		ExactlyOnce:   mf.ExactlyOnce,
		Date:          mf.Date,
		Seq:           mf.Seq,
		Slug:          mf.Slug,
		ContentSHA256: mf.ContentSHA256,
	}
}
```
