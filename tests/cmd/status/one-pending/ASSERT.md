## Expected

- Exit code **0**.
- Stdout contains a status/plan table header mentioning `MIGRATION_ID` (or `migration`).
- Stdout lists FixtureIDs[0] and an `apply` action for that pending file.
- Combined output must not claim `not implemented`.

## Side Effects

- May create `t_sql_migration_log` via EnsureTable (idempotent).
- Does not apply SQL migration files (status only).

## Errors

- stderr may be empty; no usage failure required.

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

	if len(req.FixtureIDs) != 1 {
		t.Fatalf("setup: want 1 FixtureID, got %v", req.FixtureIDs)
	}
	id := req.FixtureIDs[0]
	out := resp.Stdout
	if !strings.Contains(out, "MIGRATION_ID") && !strings.Contains(strings.ToLower(out), "migration") {
		t.Fatalf("status stdout must include table header (MIGRATION_ID):\n%s\nstderr:\n%s",
			out, resp.Stderr)
	}
	if !strings.Contains(out, id) {
		t.Fatalf("status stdout must list fixture id %q:\n%s", id, out)
	}
	// Pending file should plan as apply (status shows full plan including apply).
	foundApply := false
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, id) && strings.Contains(line, "apply") {
			foundApply = true
			break
		}
	}
	if !foundApply {
		t.Fatalf("status stdout must show action apply for %q:\n%s", id, out)
	}
	combined := out + "\n" + resp.Stderr
	if strings.Contains(strings.ToLower(combined), "not implemented") {
		t.Fatalf("status must not say not implemented:\n%s", combined)
	}
}
```
