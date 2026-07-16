## Expected

- `Run` returns no harness error and a non-nil response.
- `README.md` exists at the module root.
- README content contains each of:
  - `--dsn`
  - `--dir`

## Side Effects

- None (read-only file load).

## Errors

- Missing README or missing phrase → Assert failure with path and section label.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	requireREADMEPhrases(t, req, resp, err)
}
```
