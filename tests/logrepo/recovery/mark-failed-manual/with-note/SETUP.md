# Scenario

**Feature**: MarkFailedManual with note sets failed + note

```
MarkRunning -> MarkFailedManual(id, "rolled back by ops")
Get => status=failed, note set
```

## Preconditions

- Non-empty Note; SeedStatus=running.

## Steps

1. Set fields; expect failed + note.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "mark_failed_manual"
	req.MigrationID = leafMigrationID("mark-failed-manual")
	req.SeedStatus = "running"
	req.Note = "rolled back by ops"
	req.ExactlyOnce = false
	return nil
}
```
