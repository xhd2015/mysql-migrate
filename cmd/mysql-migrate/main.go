// Command mysql-migrate is the operator binary for MySQL migrations.
//
// It parses global flags with less-flags, applies optional env fallbacks,
// builds migrate.Config, and delegates subcommands to cli.Run.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	lessflags "github.com/xhd2015/less-flags"

	"github.com/xhd2015/mysql-migrate/cli"
	"github.com/xhd2015/mysql-migrate/migrate"
)

const programName = "mysql-migrate"

func main() {
	os.Exit(run(os.Args[1:]))
}

// run parses global flags and either prints root help or hands off to cli.Run.
// Returns a process exit code (0 success/help, 2 usage, other from cli.Run).
func run(args []string) int {
	// Binary rule: empty argv is help (exit 0), not cli.Run missing-subcommand (exit 2).
	if len(args) == 0 {
		printRootHelp(os.Stdout)
		return 0
	}

	var dsn, dir string
	remain, err := lessflags.String("--dsn", &dsn).
		String("--dir", &dir).
		HelpFunc("-h,--help", func() {
			printRootHelp(os.Stdout)
		}).
		HelpNoExit().
		StopOnFirstArg().
		Parse(args)
	if err != nil {
		if err == lessflags.ErrHelp {
			return 0
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	// Globals only (e.g. --dsn x with no subcommand) → same root help, exit 0.
	if len(remain) == 0 {
		printRootHelp(os.Stdout)
		return 0
	}

	// Env fallbacks when a flag was omitted (flag wins when both set).
	if strings.TrimSpace(dsn) == "" {
		dsn = strings.TrimSpace(os.Getenv("MIGRATE_MYSQL_DSN"))
	}
	if strings.TrimSpace(dir) == "" {
		dir = strings.TrimSpace(os.Getenv("MIGRATE_MYSQL_DIR"))
	}

	cfg := migrate.Config{
		DSN:           dsn,
		MigrationsDir: dir,
		ProgramName:   programName,
	}
	return cli.Run(cfg, remain)
}

// printRootHelp writes operator-facing usage for the binary (global flags + commands).
func printRootHelp(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [global flags] <command> [args]\n", programName)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Global flags:")
	fmt.Fprintln(w, "  --dsn string   MySQL DSN (or env MIGRATE_MYSQL_DSN)")
	fmt.Fprintln(w, "  --dir string   Migrations directory (or env MIGRATE_MYSQL_DIR)")
	fmt.Fprintln(w, "  -h, --help     Show this help")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  status       Show migration status for the configured DB")
	fmt.Fprintln(w, "  plan         Show planned apply actions for the configured DB")
	fmt.Fprintln(w, "  apply        Apply pending migrations to the configured DB")
	fmt.Fprintln(w, "  mark-done    Manually mark a migration as success")
	fmt.Fprintln(w, "  mark-failed  Manually mark a migration as failed")
	fmt.Fprintln(w, "  note         Set operator note on a migration log row")
	fmt.Fprintln(w, "  allow-retry  Clear a failed exactly-once migration for retry")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Use \"<command> -h\" for subcommand help.")
}
