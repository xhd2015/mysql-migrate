# Scenario

**Feature**: Query scans all seeded rows in order

```
# seed two values then SELECT ORDER BY id
INSERT 10, 20 -> Query -> Scanned == [10, 20]
```

## Preconditions

- Op=`query_multi`.
- SeedValues `[10, 20]`.
- Unique table.

## Steps

1. Set table, seeds, Op.
2. Assert QueryCount==2 and Scanned equals seeds.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "query_multi"
	req.Table = tableName("qmul")
	req.SeedValues = []int64{10, 20}
	t.Logf("leaf query/multi-row: table=%s seeds=%v", req.Table, req.SeedValues)
	return nil
}
```
