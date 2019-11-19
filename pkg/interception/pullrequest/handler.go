package pullrequest

import (
	"io/ioutil"
	"net/http"
)

// InterceptionHandler is an http.HandlerFunc that checks that the GitHub request
// body matches the requested fields.
//
// It recognises the following request headers:
//    X-GitHub-Event - this is provided by GitHub in its HookMechanism
//    Pullrequest-Action - this is configured on the trigger interceptor
//    Pullrequest-Repo - this is the full name of the GitHub repo e.g.
//    tektoncd/triggers.
//
// If the request matches the configuration, the body is returned, otherwise the
// an error status is returned, which causes the eventlistener not to process
// the trigger.
func InterceptionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
	}

	match, err := MatchPullRequestAction(r, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !match {
		// This means that the event listener will not continue processing
		// this trigger.
		http.Error(w, "did not match", http.StatusPreconditionFailed)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
