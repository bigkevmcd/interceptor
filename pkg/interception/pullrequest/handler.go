package pullrequest

import (
	"fmt"
	"net/http"
)

// Handler is an InterceptionFunc that checks that the GitHub request
// body matches the requested fields.
//
// It recognises the following request headers:
//    X-GitHub-Event - this is provided by GitHub in its hook-mechanism
//    Pullrequest-Action - this is configured on the trigger interceptor
//    Pullrequest-Repo - this is the full name of the GitHub repo e.g.
//    tektoncd/triggers.
//
// If the request matches the configuration, the body is returned.
func Handler(r *http.Request, body []byte) ([]byte, error) {
	match, err := MatchPullRequestAction(r, body)
	if err != nil {
		return nil, fmt.Errorf("error matching pull request: %w", err)
	}
	if !match {
		return nil, nil
	}
	return body, nil
}
