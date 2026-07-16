package migrate

import "github.com/xhd2015/mysql-migrate/migrate/sqlexec"

// Config holds connection and identity settings for migration commands.
// Callers inject an already-open sqlexec.DB (typically via sqlexec.Wrap).
// The library never opens a DSN.
type Config struct {
	DB            sqlexec.DB // required for DB subcommands; nil → usage error in cli
	MigrationsDir string
	ProgramName   string
	AppliedBy     string
}
