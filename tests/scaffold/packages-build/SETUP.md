# Scenario

**Feature**: module root builds with go build ./...

```
# from repo root (d.DOCTEST_ROOT/../..)
go build ./... -> exit 0

# stubs allowed: migrate (+ optional inventory/plan/logrepo/cli/cmd)
```

## Preconditions

- Module root is two levels above this doctest tree (`d.DOCTEST_ROOT/../..`).
- `go` is on PATH (checked by root Setup).
- Production packages under the module may be stubs; they must compile.

## Steps

1. Set `req.Mode` to `packages-build`.
2. Root `Run` executes `go build ./...` with `Dir` = module root and records
   exit code plus stdout/stderr.

## Context

- Exit criterion for P1: `go build ./...` succeeds after Config/scaffold packages exist.
- Does not assert binary behavior of `cmd/mysql-migrate` (P5+).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Mode = "packages-build"
	return nil
}
```
