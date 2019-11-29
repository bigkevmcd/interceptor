package pullrequest

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func TestHandleWithSuccess(t *testing.T) {
	event := &github.PullRequestEvent{
		Action: github.String("open"),
		Repo: &github.Repository{
			FullName: github.String("testing/testing"),
		},
		PullRequest: &github.PullRequest{
			Head: &github.PullRequestBranch{
				SHA: github.String("abc1234567"),
			},
		},
	}
	r, body := makeRequest(t, event, "pull_request", "open")

	newBody, err := Handler(r, body)

	if err != nil {
		t.Fatal(err)
	}

	shortSHA := gjson.GetBytes(newBody, "intercepted.short_sha")
	if shortSHA.Value() != "abc123" {
		t.Errorf("intercepted.commit got %s, wanted %s", shortSHA, "abc123")
	}

	// Delete the addition to simplify the return comparison.
	newBody, err = sjson.DeleteBytes(newBody, "intercepted")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(newBody, body) {
		t.Fatalf("handler got incorrect body: got %s, wanted %s", newBody, body)
	}

}
