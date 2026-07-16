# Scenario

**Feature**: apply when plan has no block — run applyables / honor skips / `--to`

```
# HasBlock=false at plan time → walk Action==apply items
plan.Build -> items apply|skip only
  -> apply loop or all-skip summary; exit 0 on success
```

## Preconditions

- No blocked / deferred items in the fixture plan.
- Leaves may seed **success** logs to exercise skip-on-reapply.
- Optional `--to <id>` only applies through that migration_id inclusive.

## Steps

1. Write pending (or success-seeded) fixtures under clear path.
2. Run `apply` [ `--to id` ].
3. Assert success exit 0 and correct applied / skipped / pending side effects.

## Context

- Mutually exclusive with `exec-fail/` and `refuse-block/`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.CloseStdin = false
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
