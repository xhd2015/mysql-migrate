# Scenario

**Feature**: Wrap returns non-nil DB that can QueryRow SELECT 1

```
# happy wrap
Wrap(sqlDB) != nil
QueryRow(SELECT 1) -> Scan -> 1
```

## Preconditions

- MySQL already ensured by parent.
- No dedicated table required.

## Steps

1. Set Op=`wrap`.
2. Assert WrapNonNil and ScanValue==1.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "wrap"
	t.Log("leaf returns-db: Wrap + SELECT 1")
	return nil
}
```
