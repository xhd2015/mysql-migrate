## Expected

- Exit code **2**.
- Combined stdout+stderr mentions note (e.g. `--note`, `missing required --note`, or `note`).

## Side Effects

- No log mutation (parse fails before logrepo).

## Errors

- Usage error about required non-empty `--note`.

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
		t.Fatalf("empty --note usage must mention note:\nstdout=%q\nstderr=%q",
			resp.Stdout, resp.Stderr)
	}
}
```
