package migrate

// Config holds connection and identity settings for migration commands.
type Config struct {
	DSN           string
	MigrationsDir string
	ProgramName   string
	AppliedBy     string
}
