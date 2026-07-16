## Expected

- Exit code **0**.
- Stdout contains `Usage`.
- Stdout lists every locked subcommand:
  `status`, `plan`, `apply`, `mark-done`, `mark-failed`, `note`, `allow-retry`.
- Stdout mentions global flags `--dsn` and `--dir`.
- Stdout does **not** mention `--local` or `--remote`.

## Side Effects

- No database access; no files written.

## Errors

- None required on stderr (may be empty).

## Exit Code

- 0

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	requireExit(t, resp, 0)
	requireRootHelpTokens(t, resp.Stdout)
}
```
