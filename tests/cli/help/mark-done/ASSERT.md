## Expected

- Exit code **0**.
- Stdout contains `Usage` and `mark-done`.
- Stdout mentions `--note` (required note for human recovery ops).
- Stdout must **not** mention `--local` or `--remote`.

## Side Effects

- None.

## Errors

- None required on stderr.

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
	if !strings.Contains(out, "mark-done") {
		t.Fatalf("mark-done help must mention mark-done:\n%s", out)
	}
	if !strings.Contains(out, "--note") {
		t.Fatalf("mark-done help must mention --note:\n%s", out)
	}
	if strings.Contains(out, "--local") || strings.Contains(out, "--remote") {
		t.Fatalf("mark-done help must not mention --local/--remote:\n%s", out)
	}
}
```
