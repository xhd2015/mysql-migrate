// Package cli implements the non-interactive MySQL migrate operator tool.
//
// Public entry point:
//
//	func Run(cfg migrate.Config, args []string) int
//
// Args do not include the program name. Exit codes:
//
//	0 — success (including help)
//	1 — business failure (HasBlock on status/plan, apply refuse/fail, recovery biz error)
//	2 — usage error, unknown command, or missing required config/flags
//
// Run never calls os.Exit.
package cli

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/xhd2015/mysql-migrate/migrate"
	"github.com/xhd2015/mysql-migrate/migrate/inventory"
	"github.com/xhd2015/mysql-migrate/migrate/logrepo"
	"github.com/xhd2015/mysql-migrate/migrate/plan"
)

// Exit codes locked by CLI tests.
const (
	ExitOK    = 0
	ExitBiz   = 1 // HasBlock, apply failure, open/DB error, recovery logrepo/biz error
	ExitUsage = 2 // usage / unknown / missing flags or config
)

// Run handles argv without the program name.
// Writes help / errors to os.Stdout and os.Stderr.
// Reads nothing required from os.Stdin (must not block if stdin is closed).
// Never calls os.Exit.
func Run(cfg migrate.Config, args []string) int {
	program := resolveProgramName(cfg.ProgramName)

	if len(args) == 0 {
		printRootUsage(os.Stdout, program)
		fmt.Fprintln(os.Stderr, "Error: missing subcommand")
		return ExitUsage
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "-h", "--help", "help":
		printRootUsage(os.Stdout, program)
		return ExitOK
	case "status":
		return runStatus(cfg, program, rest)
	case "plan":
		return runPlan(cfg, program, rest)
	case "apply":
		return runApply(cfg, program, rest)
	case "mark-done":
		return runRecovery(cfg, program, "mark-done", rest)
	case "mark-failed":
		return runRecovery(cfg, program, "mark-failed", rest)
	case "note":
		return runRecovery(cfg, program, "note", rest)
	case "allow-retry":
		return runRecovery(cfg, program, "allow-retry", rest)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown subcommand %q\n", cmd)
		printRootUsage(os.Stderr, program)
		return ExitUsage
	}
}

func resolveProgramName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "mysql-migrate"
	}
	return name
}

func printRootUsage(w io.Writer, program string) {
	fmt.Fprintf(w, "Usage: %s <command> [flags]\n", program)
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
	fmt.Fprintln(w, "Config (passed by the caller, not CLI flags):")
	fmt.Fprintln(w, "  DSN            MySQL DSN (required for DB subcommands)")
	fmt.Fprintln(w, "  MigrationsDir  Directory of *.sql migration files")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Use \"<command> -h\" for subcommand help.")
}

func wantsHelp(args []string) bool {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

func requireDSN(cfg migrate.Config) int {
	if strings.TrimSpace(cfg.DSN) == "" {
		fmt.Fprintln(os.Stderr, "Error: missing DSN on config (cfg.DSN is empty)")
		return ExitUsage
	}
	return ExitOK
}

func requireMigrationsDir(cfg migrate.Config) int {
	if strings.TrimSpace(cfg.MigrationsDir) == "" {
		fmt.Fprintln(os.Stderr, "Error: missing MigrationsDir on config (cfg.MigrationsDir is empty)")
		return ExitUsage
	}
	return ExitOK
}

// openDB opens and pings MySQL with cfg.DSN. Caller must Close.
func openDB(cfg migrate.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func runStatus(cfg migrate.Config, program string, args []string) int {
	if wantsHelp(args) {
		fmt.Fprintf(os.Stdout, "Usage: %s status\n", program)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Show migration status for the configured database.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Uses cfg.DSN and cfg.MigrationsDir from Config (no target flags).")
		return ExitOK
	}
	if hasUnknownFlags(args, nil) || len(nonFlagPositional(args)) > 0 {
		fmt.Fprintf(os.Stderr, "Error: unexpected arguments for status: %s\n", strings.Join(args, " "))
		return ExitUsage
	}
	if code := requireDSN(cfg); code != ExitOK {
		return code
	}
	if code := requireMigrationsDir(cfg); code != ExitOK {
		return code
	}
	return runStatusOrPlan(cfg, false /* planOnlyNonSkip */)
}

func runPlan(cfg migrate.Config, program string, args []string) int {
	if wantsHelp(args) {
		fmt.Fprintf(os.Stdout, "Usage: %s plan\n", program)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Show planned apply actions for the configured database.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Uses cfg.DSN and cfg.MigrationsDir from Config (no target flags).")
		return ExitOK
	}
	if hasUnknownFlags(args, nil) || len(nonFlagPositional(args)) > 0 {
		fmt.Fprintf(os.Stderr, "Error: unexpected arguments for plan: %s\n", strings.Join(args, " "))
		return ExitUsage
	}
	if code := requireDSN(cfg); code != ExitOK {
		return code
	}
	if code := requireMigrationsDir(cfg); code != ExitOK {
		return code
	}
	return runStatusOrPlan(cfg, true /* planOnlyNonSkip */)
}

// runStatusOrPlan opens the DB, builds a plan from inventory + log, prints a
// table, emits hash-mismatch warnings, and returns 1 if HasBlock else 0.
// When onlyNonSkip is true (plan subcommand), skip rows are omitted from stdout.
func runStatusOrPlan(cfg migrate.Config, onlyNonSkip bool) int {
	db, err := openDB(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: open DSN: %v\n", err)
		return ExitBiz
	}
	defer db.Close()

	if err := logrepo.EnsureTable(db); err != nil {
		fmt.Fprintf(os.Stderr, "Error: ensure migration log table: %v\n", err)
		return ExitBiz
	}

	files, err := inventory.ListDir(cfg.MigrationsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: list migrations in %s: %v\n", cfg.MigrationsDir, err)
		return ExitBiz
	}

	rows, err := logrepo.List(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: list migration log: %v\n", err)
		return ExitBiz
	}

	logs := make([]plan.LogRow, 0, len(rows))
	byID := make(map[string]logrepo.Row, len(rows))
	for _, r := range rows {
		byID[r.MigrationID] = r
		logs = append(logs, plan.LogRow{
			MigrationID:   r.MigrationID,
			Status:        r.Status,
			ExactlyOnce:   r.ExactlyOnce,
			ContentSHA256: r.ContentSHA256,
			DurationMS:    r.DurationMS,
			ErrorMessage:  r.ErrorMessage,
			Note:          r.Note,
		})
	}

	p := plan.Build(files, logs)
	printPlanTable(os.Stdout, p, byID, onlyNonSkip)
	printHashMismatchWarnings(os.Stderr, p)

	if p.HasBlock {
		return ExitBiz
	}
	return ExitOK
}

// printPlanTable writes a tab-aligned operator table. onlyNonSkip omits ActionSkip.
func printPlanTable(w io.Writer, p plan.Plan, byID map[string]logrepo.Row, onlyNonSkip bool) {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "MIGRATION_ID\tACTION\tDURATION_MS\tEXACTLY_ONCE\tNOTE")
	for _, item := range p.Items {
		if onlyNonSkip && item.Action == plan.ActionSkip {
			continue
		}
		row := byID[item.MigrationID]
		eo := "false"
		if item.ExactlyOnce {
			eo = "true"
		}
		note := strings.ReplaceAll(row.Note, "\t", " ")
		note = strings.ReplaceAll(note, "\n", " ")
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n",
			item.MigrationID,
			string(item.Action),
			row.DurationMS,
			eo,
			note,
		)
	}
	_ = tw.Flush()
}

// printHashMismatchWarnings writes one "warning:" line per HashMismatch item.
func printHashMismatchWarnings(w io.Writer, p plan.Plan) {
	for _, item := range p.Items {
		if !item.HashMismatch {
			continue
		}
		fmt.Fprintf(w, "warning: migration %s content hash mismatch (log hash differs from file)\n", item.MigrationID)
	}
}

func runApply(cfg migrate.Config, program string, args []string) int {
	if wantsHelp(args) {
		fmt.Fprintf(os.Stdout, "Usage: %s apply [--to <migration_id>]\n", program)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Apply pending migrations to the configured database.")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Flags:")
		fmt.Fprintln(os.Stdout, "  --to <id>  Apply up to and including this migration_id")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Uses cfg.DSN and cfg.MigrationsDir from Config (no target flags).")
		return ExitOK
	}
	if code := requireDSN(cfg); code != ExitOK {
		return code
	}
	if code := requireMigrationsDir(cfg); code != ExitOK {
		return code
	}
	toID, err := parseOptionalTo(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitUsage
	}
	return doApply(cfg, toID)
}

// doApply builds the same plan as status/plan, refuses when HasBlock, then
// walks Action==apply items: MarkRunning → Exec SQL → MarkSuccess|MarkFailed.
// --to stops after the named migration_id (inclusive). Progress + summary go to stdout.
func doApply(cfg migrate.Config, toID string) int {
	db, err := openDB(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: open DSN: %v\n", err)
		return ExitBiz
	}
	defer db.Close()

	if err := logrepo.EnsureTable(db); err != nil {
		fmt.Fprintf(os.Stderr, "Error: ensure migration log table: %v\n", err)
		return ExitBiz
	}

	files, err := inventory.ListDir(cfg.MigrationsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: list migrations in %s: %v\n", cfg.MigrationsDir, err)
		return ExitBiz
	}

	rows, err := logrepo.List(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: list migration log: %v\n", err)
		return ExitBiz
	}

	logs := make([]plan.LogRow, 0, len(rows))
	for _, r := range rows {
		logs = append(logs, plan.LogRow{
			MigrationID:   r.MigrationID,
			Status:        r.Status,
			ExactlyOnce:   r.ExactlyOnce,
			ContentSHA256: r.ContentSHA256,
			DurationMS:    r.DurationMS,
			ErrorMessage:  r.ErrorMessage,
			Note:          r.Note,
		})
	}

	p := plan.Build(files, logs)

	// Refuse when plan already has a block. Do not MarkRunning/Exec deferred later migrations.
	if p.HasBlock {
		blockedID := ""
		for _, item := range p.Items {
			if item.Action == plan.ActionBlocked {
				blockedID = item.MigrationID
				break
			}
		}
		if blockedID != "" {
			fmt.Fprintf(os.Stderr, "Error: plan blocked on migration %s; refuse apply (resolve with mark-done / allow-retry / fix hash)\n", blockedID)
		} else {
			fmt.Fprintln(os.Stderr, "Error: plan blocked; refuse apply")
		}
		return ExitBiz
	}

	byID := make(map[string]inventory.MigrationFile, len(files))
	for _, f := range files {
		byID[f.ID] = f
	}

	appliedBy := resolveAppliedBy(cfg.AppliedBy)
	var applied, failed, pending int

	// Count all plan-time applyables first; remaining after --to or mid-run failure is pending.
	var applyQueue []plan.PlanItem
	for _, item := range p.Items {
		if item.Action == plan.ActionApply {
			applyQueue = append(applyQueue, item)
		}
	}
	pending = len(applyQueue)

	for _, item := range applyQueue {
		f, ok := byID[item.MigrationID]
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: migration %s missing from inventory\n", item.MigrationID)
			return ExitBiz
		}

		if err := logrepo.MarkRunning(db, item.MigrationID, item.ExactlyOnce, f.ContentSHA256, appliedBy); err != nil {
			fmt.Fprintf(os.Stderr, "Error: MarkRunning %s: %v\n", item.MigrationID, err)
			return ExitBiz
		}

		sqlBytes, err := os.ReadFile(f.Path)
		if err != nil {
			_ = logrepo.MarkFailed(db, item.MigrationID, 0, err.Error())
			fmt.Fprintf(os.Stdout, "apply %s failed (read: %v)\n", item.MigrationID, err)
			failed++
			pending--
			printApplySummary(applied, failed, pending)
			return ExitBiz
		}

		start := time.Now()
		execErr := execMigrationSQL(db, string(sqlBytes))
		durMS := int(time.Since(start).Milliseconds())
		if durMS < 0 {
			durMS = 0
		}

		if execErr != nil {
			errMsg := execErr.Error()
			if markErr := logrepo.MarkFailed(db, item.MigrationID, durMS, errMsg); markErr != nil {
				fmt.Fprintf(os.Stderr, "Error: MarkFailed %s: %v\n", item.MigrationID, markErr)
			}
			fmt.Fprintf(os.Stdout, "apply %s failed (%dms): %v\n", item.MigrationID, durMS, execErr)
			failed++
			pending--
			// Stop: later Action==apply items remain pending (not executed).
			printApplySummary(applied, failed, pending)
			return ExitBiz
		}

		if err := logrepo.MarkSuccess(db, item.MigrationID, durMS); err != nil {
			fmt.Fprintf(os.Stderr, "Error: MarkSuccess %s: %v\n", item.MigrationID, err)
			return ExitBiz
		}
		fmt.Fprintf(os.Stdout, "apply %s ok (%dms)\n", item.MigrationID, durMS)
		applied++
		pending--

		// --to inclusive: stop after successfully applying this migration_id.
		if toID != "" && item.MigrationID == toID {
			break
		}
	}

	printApplySummary(applied, failed, pending)
	return ExitOK
}

func printApplySummary(applied, failed, pending int) {
	fmt.Fprintf(os.Stdout, "%d applied, %d failed, %d pending\n", applied, failed, pending)
}

// execMigrationSQL runs the full migration file body (multiStatements enabled on DSN).
func execMigrationSQL(db *sql.DB, sqlText string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	_, err := db.ExecContext(ctx, sqlText)
	return err
}

func resolveAppliedBy(cfgAppliedBy string) string {
	if s := strings.TrimSpace(cfgAppliedBy); s != "" {
		return s
	}
	return "mysql-migrate"
}

func runRecovery(cfg migrate.Config, program, cmd string, args []string) int {
	if wantsHelp(args) {
		fmt.Fprintf(os.Stdout, "Usage: %s %s <migration_id> --note \"...\"\n", program, cmd)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "%s updates the migration log for recovery / annotation.\n", cmd)
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Arguments:")
		fmt.Fprintln(os.Stdout, "  <migration_id>  Migration id (filename stem)")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Flags:")
		fmt.Fprintln(os.Stdout, "  --note string   Required operator note")
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "Uses cfg.DSN from Config (no target flags).")
		return ExitOK
	}

	// Parse id/note before DSN so usage errors (missing note/id) work offline.
	id, note, err := parseIDAndNote(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitUsage
	}

	if code := requireDSN(cfg); code != ExitOK {
		return code
	}

	db, err := openDB(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: open DSN: %v\n", err)
		return ExitBiz
	}
	defer db.Close()

	if err := logrepo.EnsureTable(db); err != nil {
		fmt.Fprintf(os.Stderr, "Error: ensure migration log table: %v\n", err)
		return ExitBiz
	}

	var opErr error
	switch cmd {
	case "mark-done":
		opErr = logrepo.MarkDone(db, id, note)
	case "mark-failed":
		opErr = logrepo.MarkFailedManual(db, id, note)
	case "note":
		opErr = logrepo.SetNote(db, id, note)
	case "allow-retry":
		opErr = logrepo.AllowRetry(db, id, note)
	default:
		fmt.Fprintf(os.Stderr, "Error: %s not implemented\n", cmd)
		return ExitUsage
	}
	if opErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", opErr)
		return ExitBiz
	}
	return ExitOK
}

func parseOptionalTo(args []string) (to string, err error) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--to":
			if i+1 >= len(args) {
				return "", fmt.Errorf("--to requires a migration_id")
			}
			to = args[i+1]
			i++
		case strings.HasPrefix(a, "--to="):
			to = strings.TrimPrefix(a, "--to=")
		case a == "-h", a == "--help":
			// ignore
		case strings.HasPrefix(a, "-"):
			return "", fmt.Errorf("unknown flag %q", a)
		default:
			return "", fmt.Errorf("unexpected argument %q", a)
		}
	}
	return to, nil
}

func parseIDAndNote(args []string) (id, note string, err error) {
	var positionals []string
	noteSet := false
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--note":
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("--note requires a value")
			}
			note = args[i+1]
			noteSet = true
			i++
		case strings.HasPrefix(a, "--note="):
			note = strings.TrimPrefix(a, "--note=")
			noteSet = true
		case a == "-h", a == "--help":
			// ignore
		case strings.HasPrefix(a, "-"):
			return "", "", fmt.Errorf("unknown flag %q", a)
		default:
			positionals = append(positionals, a)
		}
	}
	if len(positionals) == 0 {
		return "", "", fmt.Errorf("missing migration_id")
	}
	if len(positionals) > 1 {
		return "", "", fmt.Errorf("unexpected arguments: %s", strings.Join(positionals[1:], " "))
	}
	id = positionals[0]
	if !noteSet {
		return "", "", fmt.Errorf("missing required --note")
	}
	if strings.TrimSpace(note) == "" {
		return "", "", fmt.Errorf("missing required --note")
	}
	return id, note, nil
}

func nonFlagPositional(args []string) []string {
	var out []string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		out = append(out, a)
	}
	return out
}

func hasUnknownFlags(args []string, allowed map[string]bool) bool {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			continue
		}
		if a == "-h" || a == "--help" {
			continue
		}
		name := a
		if i := strings.IndexByte(a, '='); i >= 0 {
			name = a[:i]
		}
		if allowed != nil && allowed[name] {
			continue
		}
		return true
	}
	return false
}
