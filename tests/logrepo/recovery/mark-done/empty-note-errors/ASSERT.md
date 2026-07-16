## Expected

- `err != nil` with non-empty message (note required).
- Prefer message mentioning note (case-insensitive) when implementer lands API.

## Errors

- Empty note is rejected; status must not flip to success via this call.

## Side Effects

- If implementer rejects before write, row may remain `running` (not asserted strictly).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	mustErr(t, err, "MarkDone empty note")
	// Helpful error text when present.
	if msg := strings.ToLower(err.Error()); !strings.Contains(msg, "note") {
		t.Logf("warning: error does not mention note: %v", err)
	}
}
```
