# Scenario

**Feature**: `DB.Close` releases the handle so later ops fail

```
# close path
DB.Close() -> nil
DB.Exec(ctx, "SELECT 1") -> error
```

## Preconditions

- Live MySQL required (need a real connection to close).

## Steps

1. ensureMySQL.
2. Leaf sets Op=`close`.

## Context

- Single meaningful outcome under close: after-close operations error.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureMySQL(t, d)
	req.Op = "close"
	t.Log("close branch: Close then Exec must error")
	return nil
}
```
