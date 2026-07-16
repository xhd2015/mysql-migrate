# Scenario

**Feature**: migrate.Config exposes four string fields; zero value is valid

```
# zero-value Config has empty strings for all fields
var zero migrate.Config
  -> DSN, MigrationsDir, ProgramName, AppliedBy == ""

# populated Config stores the four fields as set
migrate.Config{DSN, MigrationsDir, ProgramName, AppliedBy}
  -> fields readable via Run response
```

## Preconditions

- Type `migrate.Config` exists with exported fields:
  `DSN`, `MigrationsDir`, `ProgramName`, `AppliedBy` (each `string`).
- No methods or validation required on Config for P1.

## Steps

1. Set `req.Mode` to `config-fields-exist`.
2. Provide non-empty sample values for all four fields on the request.
3. Root `Run` constructs a zero-value and a populated `Config`, recording field values.

## Context

- Sample values are arbitrary strings (no DSN parsing at P1).
- Classic RED: type missing → generated test fails to compile until implementer adds it.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "config-fields-exist"
	req.SampleDSN = "user:pass@tcp(127.0.0.1:3306)/app"
	req.SampleMigrationsDir = "migrations"
	req.SampleProgramName = "mysql-migrate"
	req.SampleAppliedBy = "ci@example.com"
	return nil
}
```
