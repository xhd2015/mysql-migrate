# Scenario

**Feature**: MarkRunning → terminal success or failed; unique migration_id upsert

```
# lifecycle upsert by migration_id
MarkRunning(id, eo, hash, by) -> status=running
MarkSuccess(id, duration_ms)   -> status=success
MarkFailed(id, duration_ms, msg) -> status=failed + error_message
# second MarkRunning same id upserts (no duplicate key error)
```

## Preconditions

- Table ensured inside Run before lifecycle ops.
- Leaves set unique `MigrationID`, hash, applied_by, duration.

## Steps

1. Child leaves set Op and field values.
2. Run seeds cleanly (DELETE id) then MarkRunning + terminal op.
3. Assert reads Get(id) fields.

## Context

- Unique constraint on migration_id is the core concurrency/retry contract.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Defaults for lifecycle leaves; children override Op and ids.
	if req.AppliedBy == "" {
		req.AppliedBy = "doctest-lifecycle"
	}
	if req.ContentSHA256 == "" {
		req.ContentSHA256 = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	}
	if req.DurationMS == 0 {
		req.DurationMS = 42
	}
	return nil
}
```
