## Expected

- Exit code **2**.
- Combined output mentions `migration_id` or `missing` (id required).

## Side Effects

- None (no DB).

## Errors

- Usage: missing migration_id.

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
	combined := strings.ToLower(resp.Stdout + "\n" + resp.Stderr)
	if !strings.Contains(combined, "migration_id") && !strings.Contains(combined, "missing") {
		t.Fatalf("missing id usage must mention migration_id or missing:\nstdout=%q\nstderr=%q",
			resp.Stdout, resp.Stderr)
	}
}
```
