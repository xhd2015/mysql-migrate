# Scenario

**Feature**: migrate.Config exposes DB + three string fields; zero value is valid

```
# zero-value Config has nil DB and empty strings
var zero migrate.Config
  -> DB == nil; MigrationsDir, ProgramName, AppliedBy == ""

# populated Config stores identity fields as set (DB may stay nil in scaffold)
migrate.Config{DB:nil, MigrationsDir, ProgramName, AppliedBy}
  -> fields readable via Run response; no DSN field on type
```

## Preconditions

- Type `migrate.Config` exists with exported fields:
  `DB` (`sqlexec.DB`), `MigrationsDir`, `ProgramName`, `AppliedBy`.
- **No** `DSN` field on Config.
- No methods or validation required on Config for scaffold.

## Steps

1. Set `req.Mode` to `config-fields-exist`.
2. Provide non-empty sample values for the three string fields on the request.
3. Root `Run` constructs a zero-value and a populated `Config`, recording field values and reflecting field presence.

## Context

- Sample values are arbitrary strings (no DSN parsing).
- Classic RED: `DB` field missing or `DSN` still present → assert fails; type
  missing → compile fails until implementer updates Config.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "config-fields-exist"
	req.SampleMigrationsDir = "migrations"
	req.SampleProgramName = "mysql-migrate"
	req.SampleAppliedBy = "ci@example.com"
	t.Log("config-fields-exist: assert DB field, no DSN")
	return nil
}
```
