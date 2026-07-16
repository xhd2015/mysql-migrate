## Expected

- Exit code **0**.
- Stdout contains `Usage` and the word `apply`.
- Stdout mentions `--to`.
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
	if !strings.Contains(out, "apply") {
		t.Fatalf("apply help must mention apply:\n%s", out)
	}
	if !strings.Contains(out, "--to") {
		t.Fatalf("apply help should mention --to:\n%s", out)
	}
	if strings.Contains(out, "--local") || strings.Contains(out, "--remote") {
		t.Fatalf("apply help must not mention --local/--remote:\n%s", out)
	}
}
```
