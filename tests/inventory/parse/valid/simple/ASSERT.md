## Expected

- `Run` succeeds (`err == nil`).
- `File.ID` is `2026-07-16-01-create-t-channel` (stem without `.sql`).
- `File.FileName` is `2026-07-16-01-create-t-channel.sql`.
- `File.Date` is `2026-07-16`.
- `File.Seq` is `1`.
- `File.ExactlyOnce` is `false`.
- `File.Slug` is `create-t-channel`.

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
	if f.ID != "2026-07-16-01-create-t-channel" {
		t.Fatalf("ID: got %q want %q", f.ID, "2026-07-16-01-create-t-channel")
	}
	if f.FileName != "2026-07-16-01-create-t-channel.sql" {
		t.Fatalf("FileName: got %q want %q", f.FileName, "2026-07-16-01-create-t-channel.sql")
	}
	if f.Date != "2026-07-16" {
		t.Fatalf("Date: got %q want %q", f.Date, "2026-07-16")
	}
	if f.Seq != 1 {
		t.Fatalf("Seq: got %d want 1", f.Seq)
	}
	if f.ExactlyOnce {
		t.Fatal("ExactlyOnce: got true want false")
	}
	if f.Slug != "create-t-channel" {
		t.Fatalf("Slug: got %q want %q", f.Slug, "create-t-channel")
	}
}
```
