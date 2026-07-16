# Scenario

**Feature**: second MarkRunning on same migration_id upserts (no duplicate error)

```
# unique(migration_id) + UPSERT semantics
MarkRunning(id, eo=false, hash1, by1)
MarkRunning(id, eo=true,  hash2, by2)  # must not fail with duplicate key
Get(id) => updated fields; COUNT(*)=1 for this id
```

## Preconditions

- Op=`unique_upsert`.
- Same MigrationID for both MarkRunning calls.
- Second call changes ExactlyOnce, ContentSHA256, AppliedBy.

## Steps

1. Set first and second MarkRunning field sets.
2. Expect no error, UniqueCount=1, row shows second-call fields, status=running.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "unique_upsert"
	req.MigrationID = leafMigrationID("unique-upsert")
	req.ExactlyOnce = false
	req.ContentSHA256 = "1111111111111111111111111111111111111111111111111111111111111111"
	req.AppliedBy = "first-run"
	req.SecondExactlyOnce = true
	req.SecondContentSHA256 = "2222222222222222222222222222222222222222222222222222222222222222"
	req.SecondAppliedBy = "retry-run"
	return nil
}
```
