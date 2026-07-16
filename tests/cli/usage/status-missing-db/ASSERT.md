## Expected

- Exit code **2**.
- Combined stdout+stderr mentions `db` or `missing` or `config` (case-insensitive).
- Must not require `--local` / `--remote` wording (those flags are gone).

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
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	requireExit(t, resp, 2)
	combined := resp.Stdout + "\n" + resp.Stderr
	lower := strings.ToLower(combined)
	if !strings.Contains(lower, "db") && !strings.Contains(lower, "missing") && !strings.Contains(lower, "config") {
		t.Fatalf("missing DB usage must mention db/missing/config:\nstdout=%q\nstderr=%q",
			resp.Stdout, resp.Stderr)
	}
	// Implementer must not reintroduce target flags as the required surface.
	if strings.Contains(combined, "--local") || strings.Contains(combined, "--remote") {
		t.Fatalf("usage error must not require --local/--remote:\n%s", combined)
	}
}
```
