# Scenario

**Feature**: README documents how to install the module / binary

```
# install section
README.md -> "go install" + module path github.com/xhd2015/mysql-migrate
```

## Preconditions

- README.md at module root (Classic RED if missing).
- Install instructions may use `go install ...@latest` or similar; required
  tokens are `go install` and the module path.

## Steps

1. Set `req.Label` to `install`.
2. Require phrases: `go install` and `github.com/xhd2015/mysql-migrate`.

## Context

- Contributors must be able to discover install from the README alone.
- Exact version suffix (`@latest` vs tag) is not locked.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Label = "install"
	req.RequiredPhrases = []string{
		"go install",
		"github.com/xhd2015/mysql-migrate",
	}
	return nil
}
```
