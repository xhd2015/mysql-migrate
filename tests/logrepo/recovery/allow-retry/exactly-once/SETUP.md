# Scenario

**Feature**: AllowRetry on EO failed row sets pending + note

```
# EO failed once — operator allows re-apply
MarkRunning(eo=true) -> MarkFailed -> AllowRetry(note)
Get => status=pending, ExactlyOnce=true, note set
```

## Preconditions

- ExactlyOnce=true on seed.
- SeedStatus=failed.
- Non-empty Note.

## Steps

1. Seed EO failed row.
2. AllowRetry with note.
3. Expect pending + note.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "allow_retry"
	req.MigrationID = leafMigrationID("allow-retry-eo")
	req.ExactlyOnce = true
	req.SeedStatus = "failed"
	req.ErrorMessage = "eo apply failed once"
	req.DurationMS = 33
	req.Note = "ops approved retry after fix"
	req.ContentSHA256 = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	return nil
}
```
