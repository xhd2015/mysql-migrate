// Package inventory discovers, parses, sorts, and hashes MySQL migration SQL files.
// It is a pure on-disk reader — no database, no apply, no CLI.
//
// Filename grammar:
//
//	YYYY-MM-DD-NN[-[EXACTLY-ONCE]]-<slug>.sql
//
// ListDir scans the top level only (not recursive), sorts by FileName ascending,
// and fills ContentSHA256. Unrelated top-level files are ignored; names that
// loosely look like migrations (YYYY-MM-DD-*.sql) but fail the full grammar
// cause ListDir to return an error.
package inventory

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// MigrationFile is structured metadata for one migration SQL file.
type MigrationFile struct {
	// ID is the migration_id (basename without .sql).
	ID string
	// FileName is the basename including .sql.
	FileName string
	// Path is the filesystem path (set by ListDir).
	Path string
	// ExactlyOnce is true when the middle [EXACTLY-ONCE] token is present.
	ExactlyOnce bool
	// Date is YYYY-MM-DD from the filename.
	Date string
	// Seq is the zero-padded day sequence NN as an int (1–99).
	Seq int
	// Slug is the kebab-case description after NN / [EXACTLY-ONCE].
	Slug string
	// ContentSHA256 is the lowercase hex SHA-256 of raw file bytes (set by ListDir).
	ContentSHA256 string
}

var (
	// exactlyOnceRe matches YYYY-MM-DD-NN-[EXACTLY-ONCE]-<slug>
	exactlyOnceRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-(\d{2})-\[EXACTLY-ONCE\]-(.+)$`)
	// simpleRe matches YYYY-MM-DD-NN-<slug>
	simpleRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})-(\d{2})-(.+)$`)
	// slugRe is kebab-case: lowercase alphanumerics with optional hyphen segments.
	slugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	// looseLookalikeRe is YYYY-MM-DD-*.sql — lookalikes that must parse or error in ListDir.
	looseLookalikeRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}-.+\.sql$`)
)

// ParseFileName parses a migration basename (with or without .sql) into metadata.
// Path and ContentSHA256 are left empty.
func ParseFileName(name string) (MigrationFile, error) {
	if name == "" {
		return MigrationFile{}, fmt.Errorf("invalid migration filename: empty name")
	}

	stem := name
	if strings.HasSuffix(name, ".sql") {
		stem = strings.TrimSuffix(name, ".sql")
	}

	var (
		date, seqStr, slug string
		exactlyOnce        bool
	)

	if m := exactlyOnceRe.FindStringSubmatch(stem); m != nil {
		date, seqStr, slug = m[1], m[2], m[3]
		exactlyOnce = true
	} else if m := simpleRe.FindStringSubmatch(stem); m != nil {
		// Incomplete / malformed [EXACTLY-ONCE] forms must not succeed as simple slugs.
		if m[3] == "[EXACTLY-ONCE]" || strings.HasPrefix(m[3], "[EXACTLY-ONCE]") {
			return MigrationFile{}, fmt.Errorf("invalid migration filename %q: missing slug after [EXACTLY-ONCE]", name)
		}
		date, seqStr, slug = m[1], m[2], m[3]
		exactlyOnce = false
	} else {
		return MigrationFile{}, fmt.Errorf("invalid migration filename %q", name)
	}

	if slug == "" || !slugRe.MatchString(slug) {
		return MigrationFile{}, fmt.Errorf("invalid migration filename %q: invalid or empty slug", name)
	}

	seq, err := strconv.Atoi(seqStr)
	if err != nil || seq < 1 || seq > 99 {
		return MigrationFile{}, fmt.Errorf("invalid migration filename %q: sequence must be 01-99", name)
	}

	return MigrationFile{
		ID:          stem,
		FileName:    stem + ".sql",
		ExactlyOnce: exactlyOnce,
		Date:        date,
		Seq:         seq,
		Slug:        slug,
	}, nil
}

// ListDir lists top-level grammar-matching .sql files in dir, sorted by FileName
// ascending, with Path and ContentSHA256 filled. Non-matching top-level junk is
// ignored. Subdirectories are not scanned. Names matching loose YYYY-MM-DD-*.sql
// that fail the full grammar return an error.
func ListDir(dir string) ([]MigrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []MigrationFile
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		name := ent.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		mf, err := ParseFileName(name)
		if err != nil {
			if looseLookalikeRe.MatchString(name) {
				return nil, fmt.Errorf("invalid migration filename %q: %w", name, err)
			}
			// Unrelated top-level .sql (no date prefix) — ignore.
			continue
		}

		path := filepath.Join(dir, name)
		hash, err := HashFile(path)
		if err != nil {
			return nil, fmt.Errorf("hash %s: %w", path, err)
		}
		mf.Path = path
		mf.ContentSHA256 = hash
		files = append(files, mf)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileName < files[j].FileName
	})
	return files, nil
}

// HashFile returns the SHA-256 of the file's raw bytes as lowercase hex (64 chars).
func HashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
