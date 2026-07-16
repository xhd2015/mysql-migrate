# mysql-migrate

**mysql-migrate** is a MySQL migration operator: it inventories SQL migration files, plans apply order, applies pending migrations, and supports operator recovery commands when a run needs a human gate.

## Architecture (core is DSN-free)

The library core (`migrate.Config`, `cli.Run`, `logrepo`, …) **never opens a DSN**.
Callers inject an already-open `sqlexec.DB` (typically via `sqlexec.Wrap(*sql.DB)`).
Only the `cmd/mysql-migrate` binary opens MySQL at the process edge from `--dsn`
(or `MIGRATE_MYSQL_DSN`), wraps it, and passes `Config.DB` into `cli.Run`.

## Install

```sh
go install github.com/xhd2015/mysql-migrate/cmd/mysql-migrate@latest
```

Module path: `github.com/xhd2015/mysql-migrate`.

## CLI usage

Global flags (before the subcommand):

| Flag | Meaning |
|------|---------|
| `--dsn` | MySQL DSN |
| `--dir` | Migrations directory |
| `-h`, `--help` | Show help |

Examples:

```sh
mysql-migrate --dsn 'user:pass@tcp(127.0.0.1:3306)/db' --dir ./migrations status
mysql-migrate --dsn 'user:pass@tcp(127.0.0.1:3306)/db' --dir ./migrations plan
mysql-migrate --dsn 'user:pass@tcp(127.0.0.1:3306)/db' --dir ./migrations apply
```

## Subcommands

| Command | Description |
|---------|-------------|
| `status` | Show migration status for the configured DB |
| `plan` | Show planned apply actions for the configured DB |
| `apply` | Apply pending migrations to the configured DB |
| `mark-done` | Manually mark a migration as success |
| `mark-failed` | Manually mark a migration as failed |
| `note` | Set operator note on a migration log row |
| `allow-retry` | Clear a failed exactly-once migration for retry |

Use `mysql-migrate <command> -h` for subcommand help.

### Migration log table

Before any DB subcommand runs, the tool ensures `t_sql_migration_log` exists
(`CREATE TABLE IF NOT EXISTS`). When the table was missing and is created, it
prints:

```text
ensured: t_sql_migration_log (created)
```

Projects should also ship an idempotent numbered migration for the same table
(so inventory/status include it). Keep that SQL in sync with the DDL embedded in
`migrate/logrepo`.

## Environment variables

When flags are omitted, the binary falls back to:

| Variable | Fallback for |
|----------|----------------|
| `MIGRATE_MYSQL_DSN` | `--dsn` |
| `MIGRATE_MYSQL_DIR` | `--dir` |

If both a flag and the corresponding env var are set, the flag wins.

## Running tests

This repo uses **doctest** design trees under `tests/` plus standard Go tests.

Validate tree structure:

```sh
doctest vet ./tests/...
```

Run the full doctest suite:

```sh
doctest test ./...
```

Run a single tree (example: docs polish):

```sh
doctest test ./tests/docs
```

Also:

```sh
go test ./...
```
