package pullrequest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bigkevmcd/interceptor/pkg/git"
	"github.com/google/go-github/github"
	"github.com/tidwall/sjson"
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
	var event github.PullRequestEvent
	err := json.Unmarshal(body, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	match, err := MatchPullRequestAction(r, body)
	if err != nil {
		return nil, fmt.Errorf("error matching pull request: %w", err)
	}
	if !match {
		return nil, nil
	}

	intercepted := map[string]interface{}{
		"short_sha": git.ShortenSHA(strValue(event.PullRequest.Head.SHA)),
	}
	body, err = sjson.SetBytes(body, "intercepted", intercepted)
	if err != nil {
		return nil, fmt.Errorf("error setting the intercepted values: %w", err)
	}

	return body, nil
}
