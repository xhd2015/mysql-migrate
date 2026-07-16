## Expected

- `Run` returns no harness error and a non-nil response.
- Package surface flags all true:
  - `HasDBType`, `HasResultType`, `HasRowsType`, `HasRowType`, `HasWrapFunc`
- Method set flags all true:
  - `DBMethodsOK` — Exec, Query, QueryRow, Close
  - `ResultMethodsOK` — LastInsertId, RowsAffected
  - `RowsMethodsOK` — Next, Scan, Close, Err
  - `RowMethodsOK` — Scan

## Side Effects

- None (pure reflection / compile-time).

## Errors

- None expected from `Run` for this mode.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireNoHarnessErr(t, err)
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	checks := []struct {
		name string
		ok   bool
	}{
		{"HasDBType", resp.HasDBType},
		{"HasResultType", resp.HasResultType},
		{"HasRowsType", resp.HasRowsType},
		{"HasRowType", resp.HasRowType},
		{"HasWrapFunc", resp.HasWrapFunc},
		{"DBMethodsOK", resp.DBMethodsOK},
		{"ResultMethodsOK", resp.ResultMethodsOK},
		{"RowsMethodsOK", resp.RowsMethodsOK},
		{"RowMethodsOK", resp.RowMethodsOK},
	}
	for _, c := range checks {
		if !c.ok {
			t.Errorf("%s = false, want true", c.name)
		}
	}
}
```
