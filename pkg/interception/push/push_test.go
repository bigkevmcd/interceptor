package push

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-github/v28/github"
)

const (
	testFullname = "testing/testing"
)

func TestMatchPushActionWithMatchingAction(t *testing.T) {
	event := makeHookBody("refs/heads/master")
	r := makeRequest(t, event, "push", "master")

	matched, err := MatchPushAction(r, event)

	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Fatal("MatchPushAction() got false, wanted true")
	}
}

func TestMatchPushActionWithUnmatchedBranch(t *testing.T) {
	event := makeHookBody("refs/heads/my-branch")
	r := makeRequest(t, event, "push", "master")

	matched, err := MatchPushAction(r, event)

	if err != nil {
		t.Fatal(err)
	}
	if matched {
		t.Fatalf("MatchPushAction() got true, wanted false")
	}
}

func TestHookKey(t *testing.T) {
	keyTests := []struct {
		event    string
		hookBody *github.PushEvent
		p        push
	}{
		{
			"push", &github.PushEvent{
				Ref: github.String("refs/heads/my-branch"),
				Repo: &github.PushEventRepository{
					FullName: github.String(testFullname),
				},
			}, push{"testing/testing", "my-branch"},
		},
	}

	for _, tt := range keyTests {
		k := pushFromHook(makeRequest(t, tt.hookBody, tt.event, "open"), tt.hookBody)
		if !k.Equal(tt.p) {
			t.Errorf("pushFromHook() got %v, wanted %v", k, tt.p)
		}
	}
}

func TestRequestKey(t *testing.T) {
	keyTests := []struct {
		eventType string
		repo      string
		ref       string
		p         push
	}{
		{"push", "testing/test", "master", push{"testing/test", "master"}},
		{"push", "testing/project", "my-branch", push{"testing/project", "my-branch"}},
	}

	for _, tt := range keyTests {
		r, _ := http.NewRequest("POST", "/testing", bytes.NewReader([]byte(``)))
		r.Header.Add("Content-Type", "application/json")
		r.Header.Add(gitHubEventHeader, tt.eventType)
		r.Header.Add(pushRepoHeader, tt.repo)
		r.Header.Add(pushRefHeader, tt.ref)

		k := pushFromRequest(r)
		if !k.Equal(tt.p) {
			t.Errorf("pushFromRequest() got %v, wanted %v", k, tt.p)
		}
	}
}

func TestRefToBranch(t *testing.T) {
	refTests := []struct {
		ref    string
		branch string
	}{
		{"refs/heads/master", "master"},
		{"refs/heads/my-branch", "my-branch"},
	}

	for _, tt := range refTests {
		if b := refToBranch(github.String(tt.ref)); b != tt.branch {
			t.Errorf("refToBranch(%s) got %s, wanted %s", tt.ref, b, tt.branch)
		}
	}
}

func makeHookBody(ref string) *github.PushEvent {
	return &github.PushEvent{
		Ref: github.String(ref),
		Repo: &github.PushEventRepository{
			FullName: github.String(testFullname),
		},
	}
}

func makeRequest(t *testing.T, event interface{}, eventType, ref string) *http.Request {
	body, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	return makeRequestWithBody(body, eventType, testFullname, ref)
}

func makeRequestWithBody(body []byte, eventType, repo, ref string) *http.Request {
	r, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(gitHubEventHeader, eventType)
	r.Header.Add(pushRefHeader, ref)
	r.Header.Add(pushRepoHeader, repo)
	return r
}
