# Scenario

**Feature**: ParseFileName accepts stem without trailing .sql

```
# same metadata whether caller passes stem or full filename
ParseFileName("2026-07-16-02-create-t-channel-participant")
  -> ID=2026-07-16-02-create-t-channel-participant
  -> FileName=2026-07-16-02-create-t-channel-participant.sql
  -> Seq=2, ExactlyOnce=false, Slug=create-t-channel-participant
```

## Preconditions

- Input has no `.sql` suffix; parser normalises `FileName` to include `.sql`.

## Steps

1. Set `FileName` to the stem form of a valid migration id.
2. Run `ParseFileName`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "parse"
	req.FileName = "2026-07-16-02-create-t-channel-participant"
	return nil
}
```
