## Expected

- `Run` succeeds.
- `Hash` and `HashB` are both valid 64-char lowercase hex.
- `Hash != HashB`.
- Digests match independent `crypto/sha256` of each body.

## Errors

- None.

```go
import (
	"crypto/sha256"
	"encoding/hex"
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
	if resp.Hash == resp.HashB {
		t.Fatalf("expected different digests, both %q", resp.Hash)
	}
	sumA := sha256.Sum256([]byte("content-A\n"))
	sumB := sha256.Sum256([]byte("content-B\n"))
	wantA := hex.EncodeToString(sumA[:])
	wantB := hex.EncodeToString(sumB[:])
	if resp.Hash != wantA {
		t.Fatalf("Hash (A): got %q want %q", resp.Hash, wantA)
	}
	if resp.HashB != wantB {
		t.Fatalf("HashB (B): got %q want %q", resp.HashB, wantB)
	}
}
```
