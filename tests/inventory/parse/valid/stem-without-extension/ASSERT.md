## Expected

- `Run` succeeds.
- `File.ID` is the stem `2026-07-16-02-create-t-channel-participant`.
- `File.FileName` is the stem with `.sql` appended.
- `File.Date` is `2026-07-16`, `Seq` is `2`, `ExactlyOnce` is false.
- `File.Slug` is `create-t-channel-participant`.

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
	if f.ID != "2026-07-16-02-create-t-channel-participant" {
		t.Fatalf("ID: got %q want %q", f.ID, "2026-07-16-02-create-t-channel-participant")
	}
	if f.FileName != "2026-07-16-02-create-t-channel-participant.sql" {
		t.Fatalf("FileName: got %q want %q", f.FileName, "2026-07-16-02-create-t-channel-participant.sql")
	}
	if f.Date != "2026-07-16" {
		t.Fatalf("Date: got %q want %q", f.Date, "2026-07-16")
	}
	if f.Seq != 2 {
		t.Fatalf("Seq: got %d want 2", f.Seq)
	}
	if f.ExactlyOnce {
		t.Fatal("ExactlyOnce: got true want false")
	}
	if f.Slug != "create-t-channel-participant" {
		t.Fatalf("Slug: got %q want %q", f.Slug, "create-t-channel-participant")
	}
}
```
