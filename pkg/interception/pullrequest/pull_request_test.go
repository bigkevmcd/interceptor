package pullrequest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-github/v28/github"
)

func TestMatchPullRequestActionWithOtherEvent(t *testing.T) {
	event := &github.PublicEvent{}
	r := makeRequest(t, event, "public", "open")

	matched, err := MatchPullRequestAction(r)

	if err != nil {
		t.Fatal(err)
	}

	if matched {
		t.Fatal("MatchPullRequestAction() got true, wanted false")
	}
}

func TestMatchPullRequestActionWithMatchingAction(t *testing.T) {
	event := makeHookBody("open")

	r := makeRequest(t, event, "pull_request", "open")

	matched, err := MatchPullRequestAction(r)

	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("MatchPullRequestAction() got false, wanted true")
	}
}

func TestMatchPullRequestActionWithDifferentAction(t *testing.T) {
	event := makeHookBody("open")
	r := makeRequest(t, event, "pull_request", "closed")

	matched, err := MatchPullRequestAction(r)

	if err != nil {
		t.Fatal(err)
	}
	if matched {
		t.Fatalf("MatchPullRequestAction() got true, wanted false")
	}
}

func TestMatchPullRequestActionInvalidJSON(t *testing.T) {
	r := makeRequestWithBody([]byte(`{test`), "pull_request", "closed")

	_, err := MatchPullRequestAction(r)
	if err == nil {
		t.Fatal("expected json parsing error, got nil")
	}

}

func TestHookKey(t *testing.T) {
	keyTests := []struct {
		event string
		body  interface{}
		key   string
	}{
		{
			"pull_request", &github.PullRequestEvent{
				Action: stringPtr("open"),
				Repo: &github.Repository{
					FullName: stringPtr("testing/testing"),
				},
			}, "pull_request:open:testing/testing",
		},
		{
			"public", &github.PublicEvent{
				Repo: &github.Repository{
					FullName: stringPtr("testing/testing"),
				},
			}, "public::testing/testing",
		},
	}

	for _, tt := range keyTests {
		r := makeRequest(t, tt.body, tt.event, "open")
		k, err := hookKey(r)
		if err != nil {
			t.Errorf("hookKey() failed: %v for case %s", err, tt.key)
		}

		if k != tt.key {
			t.Errorf("hookKey() got %s, wanted %s", k, tt.key)
		}
	}
}

func TestRequestKey(t *testing.T) {
	keyTests := []struct {
		eventType string
		repo      string
		action    string
		key       string
	}{
		{"pull_request", "testing/test", "open", "pull_request:open:testing/test"},
		{"pull_request", "testing/test", "close", "pull_request:close:testing/test"},
		{"public", "testing/test", "open", "public:open:testing/test"},
	}

	for _, tt := range keyTests {
		r, _ := http.NewRequest("POST", "/testing", bytes.NewReader([]byte(``)))
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add(gitHubEventHeader, tt.eventType)
		r.Header.Add(pullRequestRepoHeader, tt.repo)
		r.Header.Add(pullRequestActionHeader, tt.action)

		k := requestKey(r)
		if k != tt.key {
			t.Errorf("requestKey() got %s, wanted %s", k, tt.key)
		}
	}
}

func makeHookBody(action string) *github.PullRequestEvent {
	event := &github.PullRequestEvent{
		Action: stringPtr(action),
	}

	return event
}

func makeRequest(t *testing.T, event interface{}, eventType, action string) *http.Request {
	body, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	return makeRequestWithBody(body, eventType, action)
}

func makeRequestWithBody(body []byte, eventType, action string) *http.Request {
	r, _ := http.NewRequest("POST", "/testing", bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(gitHubEventHeader, eventType)
	r.Header.Add(pullRequestActionHeader, action)
	return r
}

func stringPtr(s string) *string {
	return &s
}
