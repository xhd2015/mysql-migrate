## Expected

- Exit code **0** (`HasBlock=false` on empty plan).
- Stdout does **not** list actionable migration rows for this leaf (no fixture ids).
- Combined output must **not** claim `not implemented`.

## Side Effects

- May create `t_sql_migration_log` via EnsureTable (idempotent).
- No apply of SQL migration files.

## Errors

- No usage/business failure required; stderr may be empty.

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
	requireNotStub(t, resp, "status")

	for _, id := range req.FixtureIDs {
		if strings.Contains(resp.Stdout, id) {
			t.Fatalf("empty migrations should not print fixture id %q:\n%s", id, resp.Stdout)
		}
	}

	for _, line := range strings.Split(resp.Stdout, "\n") {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		lower := strings.ToLower(trim)
		if strings.Contains(lower, "migration") || strings.Contains(lower, "action") || strings.Contains(lower, "status") {
			if !strings.Contains(trim, "2026-") {
				continue
			}
		}
		if strings.Contains(trim, "2026-") && (strings.Contains(lower, "apply") || strings.Contains(lower, "blocked") || strings.Contains(lower, "deferred")) {
			t.Fatalf("empty migrations dir should not list apply/blocked/deferred data rows:\n%s", trim)
		}
	}
}
```
