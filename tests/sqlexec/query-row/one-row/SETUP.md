# Scenario

**Feature**: QueryRow Scan returns the single seeded value

```
# happy QueryRow
INSERT 77 -> QueryRow LIMIT 1 -> Scan -> 77
```

## Preconditions

- Op=`query_row_ok`.
- SeedValues `[77]`.

## Steps

1. Set table, seed, Op.
2. Assert ScanValue == 77.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "query_row_ok"
	req.Table = tableName("qrok")
	req.SeedValues = []int64{77}
	t.Logf("leaf query-row/one-row: table=%s seed=%v", req.Table, req.SeedValues)
	return nil
}
```
