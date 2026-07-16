# Scenario

**Feature**: human recovery ops — mark-done, mark-failed-manual, set-note, allow-retry

```
# operator overrides after apply attempts (notes required for status-changing ops)
MarkDone(id, note)           -> status=success + note (note required)
MarkFailedManual(id, note)   -> status=failed + note
SetNote(id, note)            -> note only; status unchanged
AllowRetry(id, note)         -> EO failed only; status=pending + note
```

## Preconditions

- Recovery ops seed a prior row via SeedStatus (`running` | `success` | `failed`).
- MarkDone / MarkFailedManual / AllowRetry reject empty notes.
- AllowRetry requires the row's ExactlyOnce=true.

## Steps

1. Grouping sets recovery defaults; leaves set Op, SeedStatus, Note, ExactlyOnce.
2. Run seeds then calls the recovery API.
3. Assert status/note or error.

## Context

- P8 will CLI-wrap these; persistence contract is locked here.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	if req.AppliedBy == "" {
		req.AppliedBy = "doctest-recovery"
	}
	if req.ContentSHA256 == "" {
		req.ContentSHA256 = "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	}
	if req.DurationMS == 0 {
		req.DurationMS = 10
	}
	return nil
}
```
