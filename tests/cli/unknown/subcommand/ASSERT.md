## Expected

- Exit code **2** (usage).
- Stderr is non-empty and indicates an error: contains `Error` or `error`.
- Preferably mentions the unknown token or that the command is unknown.

## Side Effects

- None.

## Errors

- Usage / unknown-command error on stderr.

## Exit Code

- 2

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
	if resp.ExitCode != 2 {
		t.Fatalf("exit: got %d want 2\nstdout=%q\nstderr=%q", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	errText := resp.Stderr
	if strings.TrimSpace(errText) == "" {
		errText = resp.Stdout + resp.Stderr
	}
	lower := strings.ToLower(errText)
	if !strings.Contains(errText, "Error") && !strings.Contains(lower, "error") {
		t.Fatalf("stderr/stdout must signal Error for unknown subcommand:\nstdout=%q\nstderr=%q", resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(lower, "unknown") &&
		!strings.Contains(errText, "not-a-real-command") &&
		!strings.Contains(lower, "unrecognized") &&
		!strings.Contains(lower, "invalid") {
		t.Logf("note: message does not mention unknown/invalid/token; Error present: %q", errText)
	}
}
```
