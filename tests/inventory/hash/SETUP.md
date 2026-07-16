# Scenario

**Feature**: HashFile digests raw file bytes as SHA-256 lowercase hex

```
# content-only digest
HashFile(path) -> 64-char lowercase hex (SHA-256 of raw bytes)
```

## Preconditions

- Target path is a regular file created in a temp directory.
- Hash does not depend on filename or path.

## Steps

1. Set `req.Op = "hash"`.
2. Leaf writes one or two files and sets `Path` / `PathB`.

## Context

- Same bytes → same digest; different bytes → different digests.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "hash"
	req.Dir = newTempDir(t)
	return nil
}

// writeHashFile writes content under dir and returns absolute path.
func writeHashFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile %s: %v", path, err)
	}
	return path
}
```
