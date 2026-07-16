## Expected

- No harness error.
- `ConfigHasDB` true.
- `ConfigDBIsIface` true (`DB` field type is `sqlexec.DB`).
- `ConfigHasDSN` **false** (DSN removed from Config).
- `ConfigDBNilZero` true (zero-value DB is nil).

## Side Effects

- None.

## Errors

- None expected.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if !resp.ConfigHasDB {
		t.Fatal("migrate.Config missing field DB")
	}
	if !resp.ConfigDBIsIface {
		t.Fatal("Config.DB type must be sqlexec.DB interface")
	}
	if resp.ConfigHasDSN {
		t.Fatal("migrate.Config must not have DSN field after P1 (DB-only)")
	}
	if !resp.ConfigDBNilZero {
		t.Fatal("zero-value Config.DB must be nil")
	}
}
```
