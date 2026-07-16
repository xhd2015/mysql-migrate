## Expected

- `Run` returns no harness error and a non-nil response.
- `resp.BuildExitCode` is `0`.
- `resp.ModuleRoot` is non-empty (resolved module root path).

## Side Effects

- `go build` may write compiler artifacts under the module’s build cache only;
  no durable files under the repo are required.

## Errors

- Non-zero `BuildExitCode` is a test failure; stderr is reported for diagnosis.

## Exit Code

- `go build ./...` process exit code must be `0`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.ModuleRoot == "" {
		t.Fatal("ModuleRoot is empty")
	}
	if resp.BuildExitCode != 0 {
		t.Fatalf("go build ./... exit = %d (module root %s)\nstdout:\n%s\nstderr:\n%s",
			resp.BuildExitCode, resp.ModuleRoot, resp.BuildStdout, resp.BuildStderr)
	}
}
```
