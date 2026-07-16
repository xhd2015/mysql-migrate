# Scenario

**Feature**: MarkDone with non-empty note sets success + note

```
# seed running, then force done
MarkRunning -> MarkDone(id, "manually verified")
Get => status=success, note="manually verified"
```

## Preconditions

- Non-empty Note.
- SeedStatus=running (MarkRunning only).

## Steps

1. Set MigrationID, Note, SeedStatus=running.
2. Expect success + note.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "mark_done"
	req.MigrationID = leafMigrationID("mark-done-note")
	req.SeedStatus = "running"
	req.Note = "manually verified on staging"
	req.ExactlyOnce = false
	return nil
}
```
