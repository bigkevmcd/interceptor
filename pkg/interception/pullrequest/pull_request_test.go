package pullrequest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-github/v28/github"
)

const (
	testFullname = "testing/testing"
)

func TestMatchPullRequestActionWithOtherEvent(t *testing.T) {
	event := &github.PublicEvent{}
	r, body := makeRequest(t, event, "public", "open")

	matched, err := MatchPullRequestAction(r, body)

	if err != nil {
		t.Fatal(err)
	}

	if matched {
		t.Fatal("MatchPullRequestAction() got true, wanted false")
	}
}

func TestMatchPullRequestActionWithMatchingAction(t *testing.T) {
	event := makeHookBody("open")

	r, body := makeRequest(t, event, "pull_request", "open")

	matched, err := MatchPullRequestAction(r, body)

	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("MatchPullRequestAction() got false, wanted true")
	}
}

func TestMatchPullRequestActionWithMultipleActions(t *testing.T) {
	event := makeHookBody("synchronize")

	r, body := makeRequest(t, event, "pull_request", "open,synchronize")

	matched, err := MatchPullRequestAction(r, body)

	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("MatchPullRequestAction() got false, wanted true")
	}
}

func TestMatchPullRequestActionWithDifferentAction(t *testing.T) {
	event := makeHookBody("open")
	r, body := makeRequest(t, event, "pull_request", "closed")

	matched, err := MatchPullRequestAction(r, body)

	if err != nil {
		t.Fatal(err)
	}
	if matched {
		t.Fatalf("MatchPullRequestAction() got true, wanted false")
	}
}

func TestMatchPullRequestActionInvalidJSON(t *testing.T) {
	r, body := makeRequestWithBody([]byte(`{test`), "pull_request", testFullname, "closed")

	_, err := MatchPullRequestAction(r, body)
	if err == nil {
		t.Fatal("expected json parsing error, got nil")
	}
}

func TestExtractHookPullRequest(t *testing.T) {
	keyTests := []struct {
		event    string
		hookBody interface{}
		key      *pullRequest
	}{
		{
			"pull_request", &github.PullRequestEvent{
				Action: stringPtr("open"),
				Repo: &github.Repository{
					FullName: stringPtr(testFullname),
				},
			}, &pullRequest{"pull_request", "open", "testing/testing"},
		},
		{
			"public", &github.PublicEvent{
				Repo: &github.Repository{
					FullName: stringPtr(testFullname),
				},
			}, &pullRequest{"public", "", "testing/testing"},
		},
	}

	for _, tt := range keyTests {
		k, err := extractHookPullRequest(makeRequest(t, tt.hookBody, tt.event, "open"))
		if err != nil {
			t.Errorf("hookKey() failed: %v for case %#v", err, tt.key)
		}

		if !reflect.DeepEqual(k, tt.key) {
			t.Errorf("hookKey() got %#v, wanted %#v", k, tt.key)
		}
	}
}

func TestRequestKey(t *testing.T) {
	keyTests := []struct {
		eventType string
		repo      string
		action    string
		pr        *pullRequest
	}{
		{"pull_request", "testing/test", "open", &pullRequest{"pull_request", "open", "testing/test"}},
		{"pull_request", "testing/test", "close", &pullRequest{"pull_request", "close", "testing/test"}},
		{"public", "testing/test", "", nil},
	}

	for _, tt := range keyTests {
		r, _ := http.NewRequest("POST", "/testing", bytes.NewReader([]byte(``)))
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add(gitHubEventHeader, tt.eventType)
		r.Header.Add(pullRequestRepoHeader, tt.repo)
		r.Header.Add(pullRequestActionHeader, tt.action)

		pr := extractPullRequest(r)
		if !reflect.DeepEqual(pr, tt.pr) {
			t.Errorf("requestKey() got %#v wanted %#v", pr, tt.pr)
		}
	}
}

func makeHookBody(action string) *github.PullRequestEvent {
	event := &github.PullRequestEvent{
		Action: stringPtr(action),
		Repo: &github.Repository{
			FullName: stringPtr(testFullname),
		},
	}

	return event
}

func makeRequest(t *testing.T, event interface{}, eventType, action string) (*http.Request, []byte) {
	body, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	return makeRequestWithBody(body, eventType, testFullname, action)
}

func makeRequestWithBody(body []byte, eventType, repo, action string) (*http.Request, []byte) {
	r, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(gitHubEventHeader, eventType)
	r.Header.Add(pullRequestActionHeader, action)
	r.Header.Add(pullRequestRepoHeader, repo)
	return r, body
}

func stringPtr(s string) *string {
	return &s
}
