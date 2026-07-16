# Scenario

**Feature**: QueryRow Scan on empty match surfaces sql.ErrNoRows

```
# no-rows path
QueryRow(WHERE v = -999999).Scan(&v)
  -> errors.Is(err, sql.ErrNoRows)
```

## Preconditions

- Op=`query_row_none`.
- Empty table.

## Steps

1. Set Op and table.
2. Assert OpErrIsNoRows.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "query_row_none"
	req.Table = tableName("qrno")
	req.SeedValues = []int64{}
	t.Logf("leaf query-row/no-rows: table=%s expect ErrNoRows", req.Table)
	return nil
}
```
