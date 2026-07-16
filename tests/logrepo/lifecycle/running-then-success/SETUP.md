# Scenario

**Feature**: MarkRunning then MarkSuccess stores success, duration, hash

```
# happy path apply outcome
MarkRunning(id, eo=false, hash, by) -> MarkSuccess(id, duration_ms)
Get(id) => status=success, duration, content_sha256, applied_by
```

## Preconditions

- Op=`lifecycle_success`.
- Isolated MigrationID for this leaf.
- Non-EO row; fixed hash and applied_by.

## Steps

1. Set Op, MigrationID, ExactlyOnce=false, hash, applied_by, duration.
2. Expect success row with those fields.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "lifecycle_success"
	req.MigrationID = leafMigrationID("run-success")
	req.ExactlyOnce = false
	req.ContentSHA256 = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	req.AppliedBy = "operator-success"
	req.DurationMS = 100
	return nil
}
```
