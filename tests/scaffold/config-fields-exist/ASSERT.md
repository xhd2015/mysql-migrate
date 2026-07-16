## Expected

- `Run` returns no harness error and a non-nil response.
- Zero-value `Config`:
  - `DB` is nil (`ZeroDBIsNil`)
  - `MigrationsDir`, `ProgramName`, `AppliedBy` == `""`
- Populated `Config` echoes the sample string values set in leaf Setup.
- Reflect:
  - `HasDBField` true
  - `DBFieldIsSqlexecDB` true
  - `HasDSNField` **false**

## Side Effects

- None (pure in-memory construction).

## Errors

- None expected from `Run` for this mode.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	if !resp.ZeroDBIsNil {
		t.Fatal("zero-value Config.DB must be nil")
	}
	if resp.ZeroMigrationsDir != "" {
		t.Fatalf("zero-value Config.MigrationsDir = %q, want \"\"", resp.ZeroMigrationsDir)
	}
	if resp.ZeroProgramName != "" {
		t.Fatalf("zero-value Config.ProgramName = %q, want \"\"", resp.ZeroProgramName)
	}
	if resp.ZeroAppliedBy != "" {
		t.Fatalf("zero-value Config.AppliedBy = %q, want \"\"", resp.ZeroAppliedBy)
	}

	if resp.MigrationsDir != req.SampleMigrationsDir {
		t.Fatalf("Config.MigrationsDir = %q, want %q", resp.MigrationsDir, req.SampleMigrationsDir)
	}
	if resp.ProgramName != req.SampleProgramName {
		t.Fatalf("Config.ProgramName = %q, want %q", resp.ProgramName, req.SampleProgramName)
	}
	if resp.AppliedBy != req.SampleAppliedBy {
		t.Fatalf("Config.AppliedBy = %q, want %q", resp.AppliedBy, req.SampleAppliedBy)
	}

	if !resp.HasDBField {
		t.Fatal("migrate.Config must have field DB")
	}
	if !resp.DBFieldIsSqlexecDB {
		t.Fatal("Config.DB must be typed as sqlexec.DB")
	}
	if resp.HasDSNField {
		t.Fatal("migrate.Config must not have DSN field (DB-only)")
	}
}
```
