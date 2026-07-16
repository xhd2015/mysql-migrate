# Scenario

**Feature**: MarkRunning then MarkFailed stores failed status and error_message

```
# failed apply outcome
MarkRunning(id) -> MarkFailed(id, duration, errMsg)
Get(id) => status=failed, error_message set, duration_ms set
```

## Preconditions

- Op=`lifecycle_failed`.
- Isolated MigrationID; non-empty ErrorMessage.

## Steps

1. Set Op, MigrationID, error message, duration.
2. Expect failed row with error_message.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "lifecycle_failed"
	req.MigrationID = leafMigrationID("run-failed")
	req.ExactlyOnce = false
	req.ContentSHA256 = "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	req.AppliedBy = "operator-failed"
	req.DurationMS = 77
	req.ErrorMessage = "syntax error near line 3"
	return nil
}
```
