## Expected

- `Run` returns no harness error and a non-nil response.
- Zero-value `Config` fields are all empty strings:
  - `DSN`, `MigrationsDir`, `ProgramName`, `AppliedBy` == `""`
- Populated `Config` echoes the sample values set in leaf Setup:
  - `DSN` == `req.SampleDSN`
  - `MigrationsDir` == `req.SampleMigrationsDir`
  - `ProgramName` == `req.SampleProgramName`
  - `AppliedBy` == `req.SampleAppliedBy`

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

	if resp.ZeroDSN != "" {
		t.Fatalf("zero-value Config.DSN = %q, want \"\"", resp.ZeroDSN)
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

	if resp.DSN != req.SampleDSN {
		t.Fatalf("Config.DSN = %q, want %q", resp.DSN, req.SampleDSN)
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
}
```
