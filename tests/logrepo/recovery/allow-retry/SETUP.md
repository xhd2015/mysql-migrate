# Scenario

**Feature**: AllowRetry clears EO failed rows to pending; rejects non-EO

```
# EXACTLY-ONCE human override to re-apply
EO + failed + note -> AllowRetry -> status=pending + note
non-EO -> AllowRetry -> error
```

## Preconditions

- Op=`allow_retry` for descendants.
- Note required.
- ExactlyOnce flag on the **row** (from MarkRunning seed) gates the op.

## Steps

1. Set Op=allow_retry.
2. Leaves set ExactlyOnce, SeedStatus=failed, Note.

## Context

- Non-EO failures can re-apply via plan without allow-retry; EO must not silently clear.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "allow_retry"
	if req.SeedStatus == "" {
		req.SeedStatus = "failed"
	}
	return nil
}
```
