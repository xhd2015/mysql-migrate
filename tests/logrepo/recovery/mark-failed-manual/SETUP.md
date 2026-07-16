# Scenario

**Feature**: MarkFailedManual forces failed status with operator note

```
# recovery: force failed + note
seed row -> MarkFailedManual(id, note) -> status=failed, note set
```

## Preconditions

- Op=`mark_failed_manual`.
- Non-empty note required (same audit rule as MarkDone).

## Steps

1. Set Op=mark_failed_manual.
2. Leaf sets Note and MigrationID.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "mark_failed_manual"
	if req.SeedStatus == "" {
		req.SeedStatus = "running"
	}
	return nil
}
```
