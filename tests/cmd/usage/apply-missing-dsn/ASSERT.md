## Expected

- Exit code **2**.
- Combined stdout+stderr mentions `dsn` / `DSN` or `missing` (config incomplete).
- Must not require `--local` / `--remote` wording.

## Side Effects

- No DB open required; temp dir may be unused.

## Errors

- Usage: missing required DSN (flag and env both empty).

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
	requireExit(t, resp, 2)
	combined := resp.Stdout + "\n" + resp.Stderr
	lower := strings.ToLower(combined)
	if !strings.Contains(lower, "dsn") && !strings.Contains(lower, "missing") {
		t.Fatalf("missing DSN usage must mention dsn or missing:\nstdout=%q\nstderr=%q",
			resp.Stdout, resp.Stderr)
	}
	if strings.Contains(combined, "--local") || strings.Contains(combined, "--remote") {
		t.Fatalf("usage error must not require --local/--remote:\n%s", combined)
	}
}
```
