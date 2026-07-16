# Scenario

**Feature**: List after inserts includes every seeded migration_id

```
# three isolated ids marked running
List() contains {idA, idB, idC}
```

## Preconditions

- Primary MigrationID + two ExtraMigrationIDs, all session-prefixed.

## Steps

1. Set three unique ids under the session prefix.
2. Expect List to include all three (status running).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "list"
	req.MigrationID = leafMigrationID("list-a")
	req.ExtraMigrationIDs = []string{
		leafMigrationID("list-b"),
		leafMigrationID("list-c"),
	}
	return nil
}
```
