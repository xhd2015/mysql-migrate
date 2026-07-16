# Scenario

**Feature**: List returns migration log rows after inserts

```
# seed several MarkRunning rows, then List
MarkRunning(id1) + MarkRunning(id2) + MarkRunning(id3)
List() => includes all seeded migration_ids
```

## Preconditions

- `req.Op` is `list` for descendants.
- Unique primary MigrationID plus ExtraMigrationIDs.

## Steps

1. Set Op=list and multiple isolated ids.
2. Run seeds each id via MarkRunning then List.
3. Assert all seeded ids appear in resp.Rows.

## Context

- List returns all table rows; Assert filters by known ids (shared dev DB may hold other data).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "list"
	return nil
}
```
