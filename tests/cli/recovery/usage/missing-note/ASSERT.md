## Expected

- Exit code **2**.
- Combined output mentions `note` (required --note).

## Side Effects

- None (no DB).

## Errors

- Usage: missing required --note.

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
	if !strings.Contains(combined, "note") {
		t.Fatalf("missing --note usage must mention note:\nstdout=%q\nstderr=%q",
			resp.Stdout, resp.Stderr)
	}
}
```
