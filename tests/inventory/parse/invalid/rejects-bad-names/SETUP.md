# Scenario

**Feature**: reject a table of invalid migration basenames

```
# each bad name fails ParseFileName independently
for name in invalid_table:
  ParseFileName(name) -> error
```

## Preconditions

Invalid examples (from requirement / reference grammar):

| Name | Why invalid |
|------|-------------|
| `create-t-channel.sql` | missing date/NN |
| `2026-7-16-01-foo.sql` | unpadded month/day |
| `2026-07-16-1-foo.sql` | NN not zero-padded |
| `2026-07-16-01.sql` | missing slug |
| `2026-07-16-01-[EXACTLY-ONCE].sql` | missing slug after marker |
| `not-a-migration.txt` | wrong extension / not grammar |

## Steps

1. Set `Op` to `parse_expect_errors`.
2. Populate `InvalidNames` with the table above.
3. Run batch parse; collect per-name errors.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "parse_expect_errors"
	req.InvalidNames = []string{
		"create-t-channel.sql",
		"2026-7-16-01-foo.sql",
		"2026-07-16-1-foo.sql",
		"2026-07-16-01.sql",
		"2026-07-16-01-[EXACTLY-ONCE].sql",
		"not-a-migration.txt",
	}
	return nil
}
```
