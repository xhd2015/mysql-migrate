# Scenario

**Feature**: MarkDone forces success with a required operator note

```
# recovery: force success + note
seed row -> MarkDone(id, note) -> status=success, note set
MarkDone(id, "") -> error (note required)
```

## Preconditions

- Op=`mark_done` for descendants.
- SeedStatus typically `running` or `failed` before force-done.

## Steps

1. Set Op=mark_done.
2. Leaves set Note (or empty for error leaf) and MigrationID.

## Context

- Empty note is a hard error so audit trail is never blank.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "mark_done"
	if req.SeedStatus == "" {
		req.SeedStatus = "running"
	}
	return nil
}
```
