## Expected

- No harness error.
- Exit code **2** (usage).
- Combined stdout+stderr mentions `db` or `missing` or `config` (case-insensitive).
- Must not require `--local` / `--remote` wording.
- Must not hang (Run returns).

## Side Effects

- No DB open required.

## Errors

- Usage: missing required DB on Config.

## Exit Code

- 2

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.ExitCode != 2 {
		t.Fatalf("exit: got %d want 2\nstdout=%q\nstderr=%q",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + "\n" + resp.Stderr
	lower := strings.ToLower(combined)
	if !strings.Contains(lower, "db") &&
		!strings.Contains(lower, "missing") &&
		!strings.Contains(lower, "config") {
		t.Fatalf("nil DB usage must mention db/missing/config:\nstdout=%q\nstderr=%q",
			resp.Stdout, resp.Stderr)
	}
	if strings.Contains(combined, "--local") || strings.Contains(combined, "--remote") {
		t.Fatalf("usage error must not require --local/--remote:\n%s", combined)
	}
}
```
