package push

import (
	"log"
	"net/http"
	"regexp"

	"github.com/google/go-github/v28/github"
)

const (
	gitHubEventHeader    = "X-Github-Event"
	pushEventType        = "push"
	pushRefHeader        = "Push-Ref"
	pushExcludeRefHeader = "PushExclude-Ref"
	pushRepoHeader       = "Push-Repo"
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

	hookPush := pushFromHook(r, event)
	requestPush := pushFromRequest(r)
	log.Printf("debug: hookPush = %v, requestPush = %s", hookPush, requestPush)

	return requestMatchesHook(requestPush, hookPush), nil
}

func isPushEvent(r *http.Request) bool {
	return r.Header.Get(gitHubEventHeader) == pushEventType
}

// exclude is only set from incoming request headers in the pushFromRequest
// function, it's used to indicate a branch that should _not_ be allowed to
// match.
type push struct {
	repoName string
	ref      string
	exclude  string
}

func pushFromHook(r *http.Request, event *github.PushEvent) *push {
	return &push{repoName: repoName(event), ref: refToBranch(event.Ref)}
}

func pushFromRequest(r *http.Request) *push {
	ref := r.Header.Get(pushRefHeader)
	exclude := r.Header.Get(pushExcludeRefHeader)
	repo := r.Header.Get(pushRepoHeader)

	return &push{repoName: repo, ref: ref, exclude: exclude}
}

func repoName(e *github.PushEvent) string {
	if e.Repo == nil {
		return ""
	}
	return strValue(e.Repo.FullName)
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

func (p push) Equal(o push) bool {
	return p.repoName == o.repoName &&
		p.ref == o.ref &&
		p.exclude == o.exclude
}

// If the request (from the headers) push matches the hook push (in the body)
// then this is true.
//
// If the requested branch is empty, then return true if the repoName matches.
// This allows for matching on _all_ branches in a repo.
func requestMatchesHook(reqPush, hookPush *push) bool {
	if reqPush.Equal(*hookPush) {
		return true
	}
	if reqPush.ref == "" && reqPush.exclude == "" {
		return reqPush.repoName == hookPush.repoName
	}
	if reqPush.ref == "" && reqPush.exclude != hookPush.ref {
		return true
	}
	return false
}
