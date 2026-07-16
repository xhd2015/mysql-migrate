# Scenario

**Feature**: README documents env fallbacks MIGRATE_MYSQL_DSN and MIGRATE_MYSQL_DIR

```
# env vars section
README.md -> MIGRATE_MYSQL_DSN, MIGRATE_MYSQL_DIR
  (flag wins when both flag and env are set — prose optional)
```

## Preconditions

- README.md at module root (Classic RED if missing).
- Env names match `cmd/mysql-migrate` fallbacks.

## Steps

1. Set `req.Label` to `env-vars`.
2. Require both env var names as substrings.

## Context

- Operators must learn they can omit `--dsn` / `--dir` when env is set.
- Exact wording of "flag wins" is not locked; names are.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Label = "env-vars"
	req.RequiredPhrases = []string{
		"MIGRATE_MYSQL_DSN",
		"MIGRATE_MYSQL_DIR",
	}
	return nil
}
```
