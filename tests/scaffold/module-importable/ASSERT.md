## Expected

- `Run` returns no harness error and a non-nil response.
- `resp.PkgPath` is exactly `github.com/xhd2015/mysql-migrate/migrate`.
- `resp.TypeName` is exactly `Config`.

## Side Effects

- None.

## Errors

- None expected from `Run` for this mode.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}

	const wantPkg = "github.com/xhd2015/mysql-migrate/migrate"
	if resp.PkgPath != wantPkg {
		t.Fatalf("Config PkgPath = %q, want %q", resp.PkgPath, wantPkg)
	}
	if resp.TypeName != "Config" {
		t.Fatalf("type name = %q, want Config", resp.TypeName)
	}
}
```
