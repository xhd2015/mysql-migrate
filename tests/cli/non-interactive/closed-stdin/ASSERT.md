## Expected

- Exit code **0** (root help).
- `resp.Duration` is **strictly less than 2 seconds**.
- Stdout still contains `Usage` (proves help ran, not a timeout no-op).

## Side Effects

- None.

## Errors

- None.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
	"time"
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
	const max = 2 * time.Second
	if resp.Duration >= max {
		t.Fatalf("cli.Run hung or was too slow under closed stdin: duration=%v max=%v", resp.Duration, max)
	}
	if !strings.Contains(resp.Stdout, "Usage") {
		t.Fatalf("expected help Usage on stdout under closed stdin:\n%s", resp.Stdout)
	}
}
```
