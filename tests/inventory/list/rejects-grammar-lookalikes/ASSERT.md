## Expected

- `Run` returns a non-nil `err` (ListDir failed).
- Error relates to invalid migration filename / grammar (message not pinned to exact string).

## Errors

- ListDir fails when a top-level `YYYY-MM-DD-*.sql` lookalike fails full grammar.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil {
		t.Fatal("ListDir: expected error for grammar lookalike, got nil")
	}
	msg := err.Error()
	// Soft signal: mention the bad name or generic invalid/migration wording.
	if !strings.Contains(msg, "2026-07-16-1-bad-seq") &&
		!strings.Contains(strings.ToLower(msg), "invalid") &&
		!strings.Contains(strings.ToLower(msg), "migration") {
		t.Logf("note: error message may be refined by implementer: %q", msg)
	}
	// Still require non-nil error only as hard assert; log soft hints above.
}
```
