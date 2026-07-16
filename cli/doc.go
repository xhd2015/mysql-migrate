// Package cli is the non-interactive MySQL migrate operator library.
//
// Public entry: Run(cfg migrate.Config, args []string) int.
// Never calls os.Exit; never sql.Open. DB and MigrationsDir come from Config
// (no --local/--remote). Nil cfg.DB on DB subcommands is usage exit 2.
package cli
