## Expected

- `Run` succeeds.
- Exactly one listed file: `2026-07-16-01-keep-me.sql`.
- ID is `2026-07-16-01-keep-me`.
- Nested and junk paths never appear in `Files`.

## Side Effects

- Nested SQL and README remain on disk but are not inventoried.

## Errors

- None (junk is ignored, not an error).

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("ListDir: unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if len(resp.Files) != 1 {
		names := make([]string, len(resp.Files))
		for i, f := range resp.Files {
			names[i] = f.FileName
		}
		t.Fatalf("len(Files): got %d want 1 (%v)", len(resp.Files), names)
	}
	f := resp.Files[0]
	if f.FileName != "2026-07-16-01-keep-me.sql" {
		t.Fatalf("FileName: got %q want %q", f.FileName, "2026-07-16-01-keep-me.sql")
	}
	if f.ID != "2026-07-16-01-keep-me" {
		t.Fatalf("ID: got %q want %q", f.ID, "2026-07-16-01-keep-me")
	}
}
```
