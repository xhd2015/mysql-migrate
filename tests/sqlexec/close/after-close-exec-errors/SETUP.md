# Scenario

**Feature**: after Close, Exec returns an error

```
# post-close
Close() ok (or documented); Exec SELECT 1 -> error
```

## Preconditions

- Op=`close`.
- No table required.

## Steps

1. Reaffirm Op=`close`.
2. Assert CloseErr is nil and PostCloseExecErr is true.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "close"
	t.Log("leaf close/after-close-exec-errors")
	return nil
}
```
