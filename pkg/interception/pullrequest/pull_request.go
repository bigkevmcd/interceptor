package pullrequest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
)

const (
	gitHubEventHeader       = "X-Github-Event"
	pullRequestEventType    = "pull_request"
	pullRequestActionHeader = "Pullrequest-Action"
	pullRequestRepoHeader   = "Pullrequest-Repo"
)

// MatchPullRequestAction will match on pull-request requests if the action
// matches the action provided in the pullRequestActionHeader.
func MatchPullRequestAction(r *http.Request) (bool, error) {
	if !isPullRequestEvent(r) {
		log.Println("debug: dropping request because not a pull request event")
		return false, nil
	}

	key, err := hookKey(r)
	if err != nil {
		return false, fmt.Errorf("failed to create key: %w", err)
	}

	return requestKey(r) == key, nil
}

func isPullRequestEvent(r *http.Request) bool {
	return r.Header.Get(gitHubEventHeader) == pullRequestEventType
}

func hookKey(r *http.Request) (string, error) {
	et := r.Header.Get(gitHubEventHeader)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var event github.PullRequestEvent
	err := dec.Decode(&event)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	return fmt.Sprintf("%s:%s:%s", et, strValue(event.Action), repoName(&event)), nil
}

func repoName(e *github.PullRequestEvent) string {
	if e.Repo == nil {
		return ""
	}
	return strValue(e.Repo.FullName)
}

func requestKey(r *http.Request) string {
	et := r.Header.Get(gitHubEventHeader)
	repo := r.Header.Get(pullRequestRepoHeader)
	action := r.Header.Get(pullRequestActionHeader)

	return fmt.Sprintf("%s:%s:%s", et, action, repo)
}

func strValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
