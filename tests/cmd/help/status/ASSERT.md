## Expected

- Exit code **0**.
- Stdout contains `Usage` and the word `status`.
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
	requireExit(t, resp, 0)
	out := resp.Stdout
	if !strings.Contains(out, "Usage") {
		t.Fatalf("stdout must contain Usage:\n%s", out)
	}
	if !strings.Contains(out, "status") {
		t.Fatalf("status help must mention status:\n%s", out)
	}
	if strings.Contains(out, "--local") || strings.Contains(out, "--remote") {
		t.Fatalf("status help must not mention --local/--remote:\n%s", out)
	}
}
```
