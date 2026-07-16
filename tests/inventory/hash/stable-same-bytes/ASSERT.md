## Expected

- `Run` succeeds.
- `Hash` and `HashB` are equal, non-empty, length 64, lowercase hex.
- Digest matches `crypto/sha256` of `hello-migrate\n`.

## Errors

- None.

```go
import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("HashFile: unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.Hash == "" || resp.HashB == "" {
		t.Fatalf("empty hash: Hash=%q HashB=%q", resp.Hash, resp.HashB)
	}
	if resp.Hash != resp.HashB {
		t.Fatalf("unstable hash: %q != %q", resp.Hash, resp.HashB)
	}
	if len(resp.Hash) != 64 {
		t.Fatalf("hash length: got %d want 64", len(resp.Hash))
	}
	if !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(resp.Hash) {
		t.Fatalf("hash not lowercase hex: %q", resp.Hash)
	}
	sum := sha256.Sum256([]byte("hello-migrate\n"))
	want := hex.EncodeToString(sum[:])
	if resp.Hash != want {
		t.Fatalf("hash: got %q want %q", resp.Hash, want)
	}
}
```
