package push

import (
	"fmt"
	"net/http"
)

// Handler is an InterceptionFunc that checks that the GitHub request
// body matches the requested fields.
//
// It recognises the following request headers:
//    X-GitHub-Event - this is provided by GitHub in its hook-mechanism
//    Push-Ref - this is configured on the trigger interceptor
//    Push-Repo - this is the full name of the GitHub repo e.g.
//    tektoncd/triggers.
//
// If the request matches the configuration, the body is returned.
func Handler(r *http.Request, body []byte) ([]byte, error) {
	match, err := MatchPushAction(r, body)
	if err != nil {
		return nil, fmt.Errorf("error matching push: %w", err)
	}
	if !match {
		return nil, nil
	}
	return body, nil
}
