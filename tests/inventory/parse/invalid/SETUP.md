# Scenario

**Feature**: invalid basenames are rejected by ParseFileName

```
# bad grammar -> error, no partial MigrationFile success
ParseFileName(invalid) -> error
```

## Preconditions

- Names under this branch violate date padding, seq padding, slug, or extension rules.

## Steps

1. Switch op to `parse_expect_errors` and supply a table of bad basenames.

## Context

- Assert every listed name fails; empty error slot means unexpected success.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "parse_expect_errors"
	return nil
}
```
