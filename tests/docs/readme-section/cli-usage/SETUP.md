# Scenario

**Feature**: README documents CLI global flags --dsn and --dir

```
# CLI usage section
README.md -> documents --dsn and --dir (binary global flags)
```

## Preconditions

- README.md at module root (Classic RED if missing).
- Binary surface (from `cmd/mysql-migrate`): global `--dsn` and `--dir`.

## Steps

1. Set `req.Label` to `cli-usage`.
2. Require phrases: `--dsn` and `--dir`.

## Context

- Aligns README with the sealed binary flags in `tests/cmd`.
- Does not re-test flag parsing; only that the README teaches operators these flags.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Label = "cli-usage"
	req.RequiredPhrases = []string{
		"--dsn",
		"--dir",
	}
	return nil
}
```
