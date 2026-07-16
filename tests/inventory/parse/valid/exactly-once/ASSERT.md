## Expected

- `Run` succeeds.
- `File.ExactlyOnce` is `true`.
- `File.Slug` is `drop-legacy-tmp` (slug after the marker, not including it).
- `File.ID` is `2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp`.
- `File.Date` is `2026-07-17`.
- `File.Seq` is `1`.
- `File.FileName` is `2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp.sql`.

## Errors

- None.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("ParseFileName: unexpected error: %v", err)
	}
	if resp == nil || resp.File == nil {
		t.Fatal("expected File in response")
	}
	f := resp.File
	if !f.ExactlyOnce {
		t.Fatal("ExactlyOnce: got false want true")
	}
	if f.Slug != "drop-legacy-tmp" {
		t.Fatalf("Slug: got %q want %q", f.Slug, "drop-legacy-tmp")
	}
	if f.ID != "2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp" {
		t.Fatalf("ID: got %q want %q", f.ID, "2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp")
	}
	if f.Date != "2026-07-17" {
		t.Fatalf("Date: got %q want %q", f.Date, "2026-07-17")
	}
	if f.Seq != 1 {
		t.Fatalf("Seq: got %d want 1", f.Seq)
	}
	if f.FileName != "2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp.sql" {
		t.Fatalf("FileName: got %q want %q", f.FileName, "2026-07-17-01-[EXACTLY-ONCE]-drop-legacy-tmp.sql")
	}
}
```
