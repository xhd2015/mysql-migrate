# Scenario

**Feature**: human recovery commands — mark-done, mark-failed, note, allow-retry

```
# recovery is non-interactive; requires --note + migration_id; DSN from Config
operator -> mark-done|mark-failed|note|allow-retry <id> --note "..."
  -> logrepo.MarkDone|MarkFailedManual|SetNote|AllowRetry
  -> exit 0 on success; biz errors exit 1; usage exit 2
  -> never runs migration SQL; never prompts on stdin
  -> no --local/--remote flags
```

## Preconditions

- Subcommands: `mark-done`, `mark-failed`, `note`, `allow-retry`.
- Shared flag shape: `<migration_id> --note "..."`.
- Happy-path leaves seed log rows and may write migration fixtures for follow-up.
- Usage / empty-note leaves are offline (flag parse only).
- DB leaves use fillConfigForDB via subgroup setups.

## Steps

1. Group Setup documents recovery contract.
2. Children split by command (or shared `usage/`), then happy vs error path.
3. Multi-step leaves set `FollowUpArgs` (status or apply) after primary success.

## Context

- Sibling of `apply/` / `status/` / `plan/`.
- Significance: command → outcome (happy | biz-error | usage).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.CloseStdin = false
	if req.Args == nil {
		req.Args = []string{}
	}
	req.RecoveryNote = ""
	return nil
}
```
