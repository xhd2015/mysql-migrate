## Expected

- `Run` returns no harness error and a non-nil response.
- `README.md` exists at the module root.
- README content contains each subcommand name:
  - `status`
  - `plan`
  - `apply`
  - `mark-done`
  - `mark-failed`
  - `note`
  - `allow-retry`

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
