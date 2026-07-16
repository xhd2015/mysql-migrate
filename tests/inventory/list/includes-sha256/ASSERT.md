## Expected

- `Run` succeeds with exactly one file.
- `ContentSHA256` is non-empty, length 64, lowercase hex.
- `inventory.HashFile(Files[0].Path)` equals `ContentSHA256`.
- Independent `crypto/sha256` of fixture body matches the same digest.

## Errors

- None.

```go
import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"testing"

	"github.com/xhd2015/mysql-migrate/migrate/inventory"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("ListDir: unexpected error: %v", err)
	}
	if resp == nil || len(resp.Files) != 1 {
		t.Fatalf("expected exactly one file, got %+v", resp)
	}
	f := resp.Files[0]
	if f.ContentSHA256 == "" {
		t.Fatal("ContentSHA256 is empty")
	}
	if len(f.ContentSHA256) != 64 {
		t.Fatalf("ContentSHA256 length: got %d want 64 (%q)", len(f.ContentSHA256), f.ContentSHA256)
	}
	hexRe := regexp.MustCompile(`^[0-9a-f]{64}$`)
	if !hexRe.MatchString(f.ContentSHA256) {
		t.Fatalf("ContentSHA256 not lowercase hex: %q", f.ContentSHA256)
	}
	if f.Path == "" {
		t.Fatal("Path is empty; ListDir should set path for hashing")
	}
	got, err := inventory.HashFile(f.Path)
	if err != nil {
		t.Fatalf("HashFile(%q): %v", f.Path, err)
	}
	if got != f.ContentSHA256 {
		t.Fatalf("ContentSHA256 %q != HashFile %q", f.ContentSHA256, got)
	}
	body := req.FixtureFiles["2026-07-16-01-hashed.sql"]
	sum := sha256.Sum256([]byte(body))
	want := hex.EncodeToString(sum[:])
	if f.ContentSHA256 != want {
		t.Fatalf("ContentSHA256 %q != crypto/sha256 %q", f.ContentSHA256, want)
	}
}
```
