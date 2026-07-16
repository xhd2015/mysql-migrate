// Package logrepo owns the tool-managed t_sql_migration_log table and
// provides ensure / lifecycle / recovery helpers for MySQL migration apply.
package logrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Row maps one t_sql_migration_log row for callers and plan.LogRow conversion.
type Row struct {
	MigrationID   string
	Status        string // running | success | failed | unknown | pending
	ExactlyOnce   bool
	ContentSHA256 string
	DurationMS    int
	ErrorMessage  string
	Note          string
	AppliedBy     string
}

const createTableSQL = `
CREATE TABLE IF NOT EXISTS t_sql_migration_log (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  migration_id VARCHAR(255) NOT NULL,
  status VARCHAR(32) NOT NULL,
  exactly_once TINYINT(1) NOT NULL DEFAULT 0,
  content_sha256 VARCHAR(64) NOT NULL DEFAULT '',
  started_at DATETIME NULL,
  finished_at DATETIME NULL,
  duration_ms INT NOT NULL DEFAULT 0,
  error_message TEXT,
  note TEXT,
  applied_by VARCHAR(255) NOT NULL DEFAULT '',
  create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_migration_id (migration_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
`

const rowSelectCols = `
	migration_id, status, exactly_once, content_sha256,
	duration_ms, error_message, note, applied_by
`

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// EnsureTable creates t_sql_migration_log if missing (IF NOT EXISTS). Idempotent.
func EnsureTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	_, err := db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("ensure t_sql_migration_log: %w", err)
	}
	return nil
}

// List returns all log rows.
func List(db *sql.DB) ([]Row, error) {
	if db == nil {
		return nil, fmt.Errorf("nil db")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	q := `SELECT` + rowSelectCols + ` FROM t_sql_migration_log ORDER BY migration_id`
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list t_sql_migration_log: %w", err)
	}
	defer rows.Close()

	out := make([]Row, 0)
	for rows.Next() {
		r, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Get returns one row by migration_id. ok is false when not found.
func Get(db *sql.DB, migrationID string) (Row, bool, error) {
	if db == nil {
		return Row{}, false, fmt.Errorf("nil db")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	q := `SELECT` + rowSelectCols + ` FROM t_sql_migration_log WHERE migration_id = ?`
	row := db.QueryRowContext(ctx, q, migrationID)
	r, err := scanRow(row)
	if err == sql.ErrNoRows {
		return Row{}, false, nil
	}
	if err != nil {
		return Row{}, false, fmt.Errorf("get t_sql_migration_log %q: %w", migrationID, err)
	}
	return r, true, nil
}

// MarkRunning upserts a row as status=running for migration_id, storing
// exactly_once, content hash, and applied_by. Prefer ON DUPLICATE KEY UPDATE.
func MarkRunning(db *sql.DB, migrationID string, exactlyOnce bool, contentSHA256 string, appliedBy string) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}
	if migrationID == "" {
		return fmt.Errorf("empty migration_id")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	eo := 0
	if exactlyOnce {
		eo = 1
	}
	// Upsert: unique on migration_id. On conflict, refresh lifecycle fields
	// and leave status as running for a new attempt.
	_, err := db.ExecContext(ctx, `
		INSERT INTO t_sql_migration_log (
			migration_id, status, exactly_once, content_sha256, applied_by,
			started_at, finished_at, duration_ms, error_message
		) VALUES (?, 'running', ?, ?, ?, NOW(), NULL, 0, NULL)
		ON DUPLICATE KEY UPDATE
			status = 'running',
			exactly_once = VALUES(exactly_once),
			content_sha256 = VALUES(content_sha256),
			applied_by = VALUES(applied_by),
			started_at = NOW(),
			finished_at = NULL,
			duration_ms = 0,
			error_message = NULL
	`, migrationID, eo, contentSHA256, appliedBy)
	if err != nil {
		return fmt.Errorf("MarkRunning %q: %w", migrationID, err)
	}
	return nil
}

// MarkSuccess sets status=success and records duration_ms / finished_at.
func MarkSuccess(db *sql.DB, migrationID string, durationMS int) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	res, err := db.ExecContext(ctx, `
		UPDATE t_sql_migration_log
		SET status = 'success',
		    duration_ms = ?,
		    finished_at = NOW(),
		    error_message = NULL
		WHERE migration_id = ?
	`, durationMS, migrationID)
	if err != nil {
		return fmt.Errorf("MarkSuccess %q: %w", migrationID, err)
	}
	return requireRowsAffected(res, migrationID, "MarkSuccess")
}

// MarkFailed sets status=failed with duration and error_message.
func MarkFailed(db *sql.DB, migrationID string, durationMS int, errMsg string) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	res, err := db.ExecContext(ctx, `
		UPDATE t_sql_migration_log
		SET status = 'failed',
		    duration_ms = ?,
		    error_message = ?,
		    finished_at = NOW()
		WHERE migration_id = ?
	`, durationMS, errMsg, migrationID)
	if err != nil {
		return fmt.Errorf("MarkFailed %q: %w", migrationID, err)
	}
	return requireRowsAffected(res, migrationID, "MarkFailed")
}

// MarkDone forces status=success with a required operator note.
func MarkDone(db *sql.DB, migrationID string, note string) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}
	if err := requireNote(note); err != nil {
		return err
	}
	ctx, cancel := withTimeout()
	defer cancel()
	res, err := db.ExecContext(ctx, `
		UPDATE t_sql_migration_log
		SET status = 'success',
		    note = ?,
		    finished_at = NOW()
		WHERE migration_id = ?
	`, note, migrationID)
	if err != nil {
		return fmt.Errorf("MarkDone %q: %w", migrationID, err)
	}
	return requireRowsAffected(res, migrationID, "MarkDone")
}

// MarkFailedManual forces status=failed with a required operator note.
func MarkFailedManual(db *sql.DB, migrationID string, note string) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}
	if err := requireNote(note); err != nil {
		return err
	}
	ctx, cancel := withTimeout()
	defer cancel()
	res, err := db.ExecContext(ctx, `
		UPDATE t_sql_migration_log
		SET status = 'failed',
		    note = ?,
		    finished_at = NOW()
		WHERE migration_id = ?
	`, note, migrationID)
	if err != nil {
		return fmt.Errorf("MarkFailedManual %q: %w", migrationID, err)
	}
	return requireRowsAffected(res, migrationID, "MarkFailedManual")
}

// SetNote updates note without changing status.
func SetNote(db *sql.DB, migrationID string, note string) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}
	ctx, cancel := withTimeout()
	defer cancel()
	res, err := db.ExecContext(ctx, `
		UPDATE t_sql_migration_log
		SET note = ?
		WHERE migration_id = ?
	`, note, migrationID)
	if err != nil {
		return fmt.Errorf("SetNote %q: %w", migrationID, err)
	}
	return requireRowsAffected(res, migrationID, "SetNote")
}

// AllowRetry sets status=pending + note for an exactly-once failed row so plan
// can re-apply. Non-exactly-once rows return an error. Note is required.
func AllowRetry(db *sql.DB, migrationID string, note string) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}
	if err := requireNote(note); err != nil {
		return err
	}
	row, ok, err := Get(db, migrationID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("AllowRetry: migration_id %q not found", migrationID)
	}
	if !row.ExactlyOnce {
		return fmt.Errorf("AllowRetry: migration_id %q is not exactly-once (EO); cannot clear for retry", migrationID)
	}
	ctx, cancel := withTimeout()
	defer cancel()
	res, err := db.ExecContext(ctx, `
		UPDATE t_sql_migration_log
		SET status = 'pending',
		    note = ?
		WHERE migration_id = ?
	`, note, migrationID)
	if err != nil {
		return fmt.Errorf("AllowRetry %q: %w", migrationID, err)
	}
	return requireRowsAffected(res, migrationID, "AllowRetry")
}

func requireNote(note string) error {
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("note is required and must be non-empty")
	}
	return nil
}

func requireRowsAffected(res sql.Result, migrationID, op string) error {
	n, err := res.RowsAffected()
	if err != nil {
		// Some drivers may not report rows affected; treat as success.
		return nil
	}
	if n == 0 {
		return fmt.Errorf("%s: migration_id %q not found", op, migrationID)
	}
	return nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanRow(s scannable) (Row, error) {
	var (
		r      Row
		eo     int
		errMsg sql.NullString
		note   sql.NullString
	)
	err := s.Scan(
		&r.MigrationID,
		&r.Status,
		&eo,
		&r.ContentSHA256,
		&r.DurationMS,
		&errMsg,
		&note,
		&r.AppliedBy,
	)
	if err != nil {
		return Row{}, err
	}
	r.ExactlyOnce = eo != 0
	if errMsg.Valid {
		r.ErrorMessage = errMsg.String
	}
	if note.Valid {
		r.Note = note.String
	}
	return r, nil
}
