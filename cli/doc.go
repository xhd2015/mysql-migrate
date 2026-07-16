// Package cli implements the non-interactive mysql-migrate operator CLI library.
//
// Entry point: Run(cfg migrate.Config, args []string) int
// Never calls os.Exit; DSN and MigrationsDir come from Config (no --local/--remote).
package cli
