package push

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func TestHandleWithSuccess(t *testing.T) {
	event := &github.PushEvent{
		Ref: stringPtr("refs/heads/master"),
		Repo: &github.PushEventRepository{
			FullName: stringPtr("testing/testing"),
		},
		Commits: []github.PushEventCommit{
			{ID: stringPtr("abc123")},
			{ID: stringPtr("cde23456")},
		},
	}
	r := makeRequest(t, event, "push", "master")
	body := mustMarshal(t, event)
	newBody, err := Handler(r, body)

	if err != nil {
		t.Fatal(err)
	}

	interceptedRef := gjson.GetBytes(newBody, "intercepted.ref")
	if interceptedRef.Value() != "master" {
		t.Errorf("intercepted.ref got %s, wanted %s", interceptedRef, "master")
	}
	interceptedCommit := gjson.GetBytes(newBody, "intercepted.last_commit")
	if interceptedCommit.Value() != "abc123" {
		t.Errorf("intercepted.commit got %s, wanted %s", interceptedCommit, "abc123")
	}

	// Delete the addition to simplify the return comparison.
	returnBody, err := sjson.DeleteBytes(newBody, "intercepted")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(returnBody, body) {
		t.Fatalf("handler got incorrect body: got %s, wanted %s", newBody, body)
	}

}

func TestExtractSecondLastCommit(t *testing.T) {
	commitID := "abc12345"
	commitTests := []struct {
		commits   []github.PushEventCommit
		extracted string
	}{
		{
			[]github.PushEventCommit{
				{ID: stringPtr(commitID)},
				{ID: stringPtr("cde23456")},
			},
			commitID,
		},
		{
			[]github.PushEventCommit{
				{ID: stringPtr(commitID)},
			},
			commitID,
		},
	}

	for _, tt := range commitTests {
		event := &github.PushEvent{
			Ref: stringPtr("refs/heads/master"),
			Repo: &github.PushEventRepository{
				FullName: stringPtr("testing/testing"),
			},
			Commits: tt.commits,
		}

		result := secondLastCommit(event)

		if result != tt.extracted {
			t.Errorf("secondLastCommit got %s, wanted %s", result, tt.extracted)
		}
	}

}

func mustMarshal(t *testing.T, e interface{}) []byte {
	body, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	return body
}
