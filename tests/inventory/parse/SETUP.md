# Scenario

**Feature**: ParseFileName turns a migration basename into structured metadata

```
# parse basename (with or without .sql) into MigrationFile fields
ParseFileName(name) -> {ID, FileName, Date, Seq, ExactlyOnce, Slug} | error
```

## Preconditions

- Input is a basename only (not a full path).
- Extension may be present (`.sql`) or omitted (stem form).
- `Path` and `ContentSHA256` are not required from ParseFileName.

## Steps

1. Set `req.Op` to `parse` (invalid branch may switch to `parse_expect_errors`).
2. Child nodes set `FileName` or `InvalidNames`.

## Context

- Valid grammar:
  `YYYY-MM-DD-NN[-[EXACTLY-ONCE]]-<slug>[.sql]`
- Invalid names must return a non-nil error (message content not tightly fixed).

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "parse"
	return nil
}
```
