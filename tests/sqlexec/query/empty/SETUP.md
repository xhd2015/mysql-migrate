# Scenario

**Feature**: Query with no matching rows is success with empty cursor

```
# empty result set
Query(WHERE v = -999999) -> Next never true; Err nil; Close ok
```

## Preconditions

- Op=`query_empty`.
- Empty table (no seed values).

## Steps

1. Set Op and table.
2. Assert QueryEmpty true and QueryCount 0.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "query_empty"
	req.Table = tableName("qemp")
	req.SeedValues = []int64{} // empty table
	t.Logf("leaf query/empty: table=%s", req.Table)
	return nil
}
```
