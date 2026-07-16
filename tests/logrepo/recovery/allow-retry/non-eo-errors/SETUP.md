# Scenario

**Feature**: AllowRetry on non-EO row returns error

```
# non-EO must not use allow-retry path
MarkRunning(eo=false) -> MarkFailed -> AllowRetry(note) -> error
```

## Preconditions

- ExactlyOnce=false.
- SeedStatus=failed.
- Note non-empty (failure is EO gate, not note).

## Steps

1. Seed non-EO failed row.
2. Expect AllowRetry error.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "allow_retry"
	req.MigrationID = leafMigrationID("allow-retry-non-eo")
	req.ExactlyOnce = false
	req.SeedStatus = "failed"
	req.ErrorMessage = "non-eo failed"
	req.DurationMS = 11
	req.Note = "should not clear non-eo via allow-retry"
	return nil
}
```
