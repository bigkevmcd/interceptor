package push

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/google/go-github/v28/github"
)

const (
	gitHubEventHeader = "X-Github-Event"
	pushEventType     = "push"
	pushRefHeader     = "Push-Ref"
	pushRepoHeader    = "Push-Repo"
)

var branchRE = regexp.MustCompile("^refs/heads/")

// MatchPushAction will match on push notifications, if the ref for the
// commit matches the branch provided in the pushRefHeader and the Push-Repo
// matches the repository.full_name in the body.
func MatchPushAction(r *http.Request, event *github.PushEvent) (bool, error) {
	if !isPushEvent(r) {
		log.Println("debug: dropping request because not a push event")
		return false, nil
	}

	hookKey := keyFromHook(r, event)
	log.Printf("debug: hookKey = %s, requestKey = %s", hookKey, requestKey(r))

	return requestKey(r) == hookKey, nil
}

func isPushEvent(r *http.Request) bool {
	return r.Header.Get(gitHubEventHeader) == pushEventType
}

func keyFromHook(r *http.Request, event *github.PushEvent) string {
	et := r.Header.Get(gitHubEventHeader)
	return fmt.Sprintf("%s:%s:%s", et, repoName(event), refToBranch(event.Ref))
}

func repoName(e *github.PushEvent) string {
	if e.Repo == nil {
		return ""
	}
	return strValue(e.Repo.FullName)
}

func requestKey(r *http.Request) string {
	et := r.Header.Get(gitHubEventHeader)
	ref := r.Header.Get(pushRefHeader)
	repo := r.Header.Get(pushRepoHeader)

	return fmt.Sprintf("%s:%s:%s", et, repo, ref)
}

func strValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func refToBranch(s *string) string {
	return branchRE.ReplaceAllString(strValue(s), "")
}
