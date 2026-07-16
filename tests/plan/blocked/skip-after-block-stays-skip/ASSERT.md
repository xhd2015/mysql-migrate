## Expected

- `Run` succeeds.
- Actions in order: **blocked**, **skip**, **deferred**.
- Middle item must **not** be rewritten to deferred.
- `HasBlock` is true.

## Errors

- None.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("plan.Build: unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if !resp.HasBlock {
		t.Fatal("HasBlock: got false want true")
	}
	if len(resp.Items) != 3 {
		t.Fatalf("len(Items): got %d want 3", len(resp.Items))
	}
	want := []string{"blocked", "skip", "deferred"}
	for i, a := range want {
		if resp.Items[i].Action != a {
			t.Fatalf("Items[%d].Action: got %q want %q", i, resp.Items[i].Action, a)
		}
	}
	if resp.Items[1].LogStatus != "success" {
		t.Fatalf("Items[1].LogStatus: got %q want success", resp.Items[1].LogStatus)
	}
	if resp.Items[1].HashMismatch {
		t.Fatal("Items[1].HashMismatch: want false")
	}
}
```
