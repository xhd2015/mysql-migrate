# Scenario

**Feature**: SetNote changes note without changing status

```
# annotate without status flip
seed success -> SetNote(id, "post-apply note")
Get => status still success, note updated
```

## Preconditions

- Op=`set_note`.
- SeedStatus=`success` so status must remain success.
- Non-empty Note for the update.

## Steps

1. Seed success row, then SetNote.
2. Expect status unchanged and note updated.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "set_note"
	req.MigrationID = leafMigrationID("set-note-only")
	req.SeedStatus = "success"
	req.DurationMS = 5
	req.Note = "post-apply annotation"
	req.ExactlyOnce = false
	req.ContentSHA256 = "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	return nil
}
```
