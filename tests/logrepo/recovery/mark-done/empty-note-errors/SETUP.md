# Scenario

**Feature**: MarkDone with empty note returns error

```
# note required for audit
MarkRunning -> MarkDone(id, "") -> error
```

## Preconditions

- Note is empty string.
- Row seeded so failure is about the note, not missing row.

## Steps

1. Set Note="".
2. Expect Run error (non-nil).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "mark_done"
	req.MigrationID = leafMigrationID("mark-done-empty")
	req.SeedStatus = "running"
	req.Note = ""
	req.ExactlyOnce = false
	return nil
}
```
