## Expected

- `err != nil` with non-empty message.
- Prefer message mentioning exactly-once / exactly_once / EO when implementer lands API.

## Errors

- Non-EO rows cannot be cleared via AllowRetry.

## Side Effects

- Row should remain failed if implementer rejects before write (not strictly asserted).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	mustErr(t, err, "AllowRetry non-EO")
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "exactly") && !strings.Contains(msg, "eo") {
		t.Logf("warning: error does not mention exactly-once/EO: %v", err)
	}
}
```
