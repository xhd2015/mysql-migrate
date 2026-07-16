## Expected

- Exit code **0**.
- Stdout contains `Usage`.
- Stdout lists every locked subcommand:
  `status`, `plan`, `apply`, `mark-done`, `mark-failed`, `note`, `allow-retry`.
- Stdout mentions ProgramName (`mysql-migrate`).

## Side Effects

- No database access; no files written.

## Errors

- None required on stderr (may be empty).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit: got %d want 0\nstdout=%q\nstderr=%q", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	out := resp.Stdout
	if !strings.Contains(out, "Usage") {
		t.Fatalf("stdout must contain Usage:\n%s", out)
	}
	if req.ProgramName != "" && !strings.Contains(out, req.ProgramName) {
		t.Fatalf("root help should mention ProgramName %q:\n%s", req.ProgramName, out)
	}
	for _, cmd := range []string{
		"status",
		"plan",
		"apply",
		"mark-done",
		"mark-failed",
		"note",
		"allow-retry",
	} {
		if !strings.Contains(out, cmd) {
			t.Fatalf("root help must list %q:\n%s", cmd, out)
		}
	}
	// Must not document removed target flags.
	if strings.Contains(out, "--local") || strings.Contains(out, "--remote") {
		t.Fatalf("root help must not mention --local/--remote:\n%s", out)
	}
}
```
