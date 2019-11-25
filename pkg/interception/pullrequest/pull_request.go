package pullrequest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
func MatchPullRequestAction(r *http.Request, body []byte) (bool, error) {
	if !isPullRequestEvent(r) {
		log.Println("debug: dropping request because not a pull request event")
		return false, nil
	}

	hookPullRequest, err := extractHookPullRequest(r, body)
	if err != nil {
		return false, fmt.Errorf("failed to create key: %w", err)
	}

	wantedPullRequest := extractPullRequest(r)
	if wantedPullRequest == nil {
		return false, nil
	}
	log.Printf("debug: hook = %s, wanted = %s", hookPullRequest, wantedPullRequest)
	return matchHookAndRequest(hookPullRequest, wantedPullRequest), nil
}

func isPullRequestEvent(r *http.Request) bool {
	return r.Header.Get(gitHubEventHeader) == pullRequestEventType
}

type pullRequest struct {
	eventType string
	action    string
	repoName  string
}

func repoName(e *github.PullRequestEvent) string {
	if e.Repo == nil {
		return ""
	}
	return strValue(e.Repo.FullName)
}

func extractPullRequest(r *http.Request) *pullRequest {
	et := r.Header.Get(gitHubEventHeader)
	repo := r.Header.Get(pullRequestRepoHeader)
	action := r.Header.Get(pullRequestActionHeader)
	if et != pullRequestEventType {
		return nil
	}

	return &pullRequest{
		eventType: et,
		action:    action,
		repoName:  repo,
	}
}

func extractHookPullRequest(r *http.Request, body []byte) (*pullRequest, error) {
	et := r.Header.Get(gitHubEventHeader)
	var event github.PullRequestEvent
	err := json.Unmarshal(body, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	return &pullRequest{
		eventType: et,
		action:    strValue(event.Action),
		repoName:  repoName(&event),
	}, nil
}

func strValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func matchAction(header, action string) bool {
	headerActions := strings.Split(header, ",")
	for _, i := range headerActions {
		if strings.TrimSpace(i) == action {
			return true
		}
	}
	return false
}

func matchHookAndRequest(hook, req *pullRequest) bool {
	if hook.eventType != req.eventType {
		return false
	}
	if hook.repoName != req.repoName {
		return false
	}
	return matchAction(req.action, hook.action)
}
