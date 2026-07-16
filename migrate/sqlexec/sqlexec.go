// Package sqlexec is a thin context-first facade over database/sql execution.
// Callers open connections externally and inject handles via Wrap(*sql.DB).
// The migrate library never takes a DSN or calls sql.Open.
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
func Wrap(db *sql.DB) DB {
	return &dbWrap{db: db}
}

type dbWrap struct {
	db *sql.DB
}

func (w *dbWrap) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	return w.db.ExecContext(ctx, query, args...)
}

func (w *dbWrap) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	return w.db.QueryContext(ctx, query, args...)
}

func (w *dbWrap) QueryRow(ctx context.Context, query string, args ...any) Row {
	return w.db.QueryRowContext(ctx, query, args...)
}

func (w *dbWrap) Close() error {
	return w.db.Close()
}
