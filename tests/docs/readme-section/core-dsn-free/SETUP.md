# Scenario

**Feature**: README documents core Config is DSN-free (edge opens DSN)

```
# architecture note for library vs binary
README.md -> core/library never opens a DSN; binary edge opens and injects DB
```

## Preconditions

- README.md at module root (Classic RED if missing phrase).
- Aligns with `migrate.Config` having `DB sqlexec.DB` and no DSN field.
- Binary (`cmd/mysql-migrate`) owns `sql.Open` + `sqlexec.Wrap`.

## Steps

1. Set `req.Label` to `core-dsn-free`.
2. Require phrases locking the DSN-free core contract.

## Context

- P2 docs backfill: operators/contributors see that DSN is an edge concern only.
- Does not re-test Config reflect fields (`tests/scaffold`, `tests/sqlexec`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Label = "core-dsn-free"
	req.RequiredPhrases = []string{
		"never opens a DSN",
		"sqlexec",
	}
	return nil
}
```
