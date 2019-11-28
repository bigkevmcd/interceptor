package push

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/v28/github"
	"github.com/tidwall/sjson"
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
// If the request matches the configuration, the body is returned, with an
// additional key added to the body: "intercepted.ref" which will be ths
// shortened version of the ref extracting just the last part (the branch).
func Handler(r *http.Request, body []byte) ([]byte, error) {
	var event github.PushEvent
	err := json.Unmarshal(body, &event)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	match, err := MatchPushAction(r, &event)
	if err != nil {
		return nil, fmt.Errorf("error matching push: %w", err)
	}
	if !match {
		return nil, nil
	}

	intercepted := map[string]interface{}{
		"ref":       refToBranch(event.Ref),
		"short_sha": shortSHA(event.HeadCommit.ID),
	}
	updatedBody, err := sjson.SetBytes(body, "intercepted", intercepted)
	if err != nil {
		return nil, fmt.Errorf("error setting the intercepted values: %w", err)
	}
	return updatedBody, nil
}

func shortSHA(s *string) string {
	return strValue(s)[:6]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
