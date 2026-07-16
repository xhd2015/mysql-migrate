## Expected

- Exit code **0**.
- Stdout contains the exact line:
  `ensured: t_sql_migration_log (created)`
- Proves binary edge open → Wrap → cli ensure path prints on first create.

## Side Effects

- Recreates `t_sql_migration_log` in the harness schema.
- Empty migrations dir → no inventory rows required.

## Errors

- stderr may be empty.

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

	const want = "ensured: t_sql_migration_log (created)"
	if !strings.Contains(resp.Stdout, want) {
		t.Fatalf("want ensure-created line %q on stdout via Wrap path\nstdout:\n%s\nstderr:\n%s",
			want, resp.Stdout, resp.Stderr)
	}
}
```
