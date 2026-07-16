## Expected

- Exit code **0** (not usage **2**).
- Stdout is root help: `Usage`, all subcommands, `--dsn`, `--dir`.
- Proves main owns empty-args help instead of forwarding empty remain to
  `cli.Run` (which would print missing-subcommand and exit 2).

## Side Effects

- None.

## Errors

- stderr may be empty; must not require a usage Error for empty args.

## Exit Code

- 0

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode != 0 {
		t.Fatalf("empty args must exit 0 (help), got %d\nstdout=%q\nstderr=%q",
			resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	requireRootHelpTokens(t, resp.Stdout)
}
```
