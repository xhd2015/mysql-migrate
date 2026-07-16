# Scenario

**Feature**: README section leaves lock required documentation phrases

```
# each leaf under readme-section asserts one topic's phrases
README.md -> section topic (purpose | install | cli | subcommands | env | doctests)
  -> RequiredPhrases must all appear as substrings
```

## Preconditions

- Module root was validated by root Setup.
- Every descendant leaf sets non-empty `req.Label` and `req.RequiredPhrases`.
- Assert helpers live here so each leaf only names its phrases.

## Steps

1. Grouping setup only confirms the request is non-nil (leaves fill phrases).
2. After leaf Setup, root `Run` loads README; leaf Assert uses
   `requireREADMEPhrases`.

## Context

- Grouping node (no ASSERT.md): MECE split is by README section topic.
- Helper `requireREADMEPhrases` fails clearly when the file is missing or when
  a required substring is absent (shows path + label + missing phrase).

```go
import (
	"fmt"
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	if req == nil {
		return fmt.Errorf("nil request")
	}
	return nil
}

// requireREADMEPhrases fails if README is missing or any required phrase is absent.
func requireREADMEPhrases(t *testing.T, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if req == nil {
		t.Fatal("nil request")
	}
	if len(req.RequiredPhrases) == 0 {
		t.Fatal("RequiredPhrases is empty; leaf Setup must set phrases")
	}

	label := req.Label
	if label == "" {
		label = "readme"
	}

	if !resp.Exists {
		t.Fatalf("README.md missing at %s (section %q); implementer must add root README with required phrases",
			resp.READMEPath, label)
	}
	if resp.Content == "" {
		t.Fatalf("README.md at %s is empty (section %q)", resp.READMEPath, label)
	}

	for _, phrase := range req.RequiredPhrases {
		if phrase == "" {
			t.Fatal("empty RequiredPhrases entry")
		}
		if !strings.Contains(resp.Content, phrase) {
			t.Fatalf("README.md (%s) section %q must contain %q\npath: %s\n--- README begin ---\n%s\n--- README end ---",
				resp.READMEPath, label, phrase, resp.READMEPath, resp.Content)
		}
	}
}
```
