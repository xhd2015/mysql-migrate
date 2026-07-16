## Expected

- `Run` succeeds.
- Exactly three files listed.
- `Files[i].FileName` in ascending order:
  1. `2026-07-16-01-z.sql`
  2. `2026-07-16-02-a.sql`
  3. `2026-07-17-01-b.sql`
- IDs match stems without `.sql`.

## Errors

- None.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("ListDir: unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	want := []string{
		"2026-07-16-01-z.sql",
		"2026-07-16-02-a.sql",
		"2026-07-17-01-b.sql",
	}
	if len(resp.Files) != len(want) {
		t.Fatalf("len(Files): got %d want %d (%v)", len(resp.Files), len(want), fileNames(resp.Files))
	}
	for i, name := range want {
		if resp.Files[i].FileName != name {
			t.Fatalf("Files[%d].FileName: got %q want %q (all=%v)", i, resp.Files[i].FileName, name, fileNames(resp.Files))
		}
		stem := name[:len(name)-len(".sql")]
		if resp.Files[i].ID != stem {
			t.Fatalf("Files[%d].ID: got %q want %q", i, resp.Files[i].ID, stem)
		}
	}
}

func fileNames(files []FileView) []string {
	out := make([]string, len(files))
	for i, f := range files {
		out[i] = f.FileName
	}
	return out
}
```
