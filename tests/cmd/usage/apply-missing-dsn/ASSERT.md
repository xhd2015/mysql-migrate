## Expected

- Exit code **2**.
- Combined stdout+stderr mentions `dsn` / `missing` / `db` (nil cfg.DB usage path).
- Must not require `--local` / `--remote` wording.

## Side Effects

- No DB open required; temp dir may be unused.

## Errors

- Usage: missing DSN at edge leaves `cfg.DB` nil; cli requireDB fails.

## Exit Code

- 2

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	requireMissingDSNUsage(t, resp)
}
```
