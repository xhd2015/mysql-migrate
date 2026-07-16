# Scenario

**Feature**: ListDir fills ContentSHA256 matching HashFile for each entry

```
# listed hash equals independent HashFile of the same path
write 2026-07-16-01-hashed.sql with known body
ListDir -> files[0].ContentSHA256
HashFile(path) -> same hex
```

## Preconditions

- One valid migration with known non-empty SQL body.

## Steps

1. Write a single valid migration file with fixed content.
2. ListDir; capture `ContentSHA256` and path.
3. Assert leaf also re-hashes via `inventory.HashFile` to verify equality with
   `ContentSHA256`, and compares against independent `crypto/sha256`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "list"
	req.FixtureFiles = map[string]string{
		"2026-07-16-01-hashed.sql": "SELECT 1;\n-- inventory hash fixture\n",
	}
	return nil
}
```
