package interception

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bigkevmcd/interceptor/pkg/interception/pullrequest"
	"github.com/bigkevmcd/interceptor/pkg/interception/push"
)

const (
	gitHubEventHeader = "X-Github-Event"
)

// eventHandlerMap is a mapping from GitHub hook events to handlers.
var eventHandlerMap = map[string]InterceptionFunc{
	"pull_request": pullrequest.Handler,
	"push":         push.Handler,
}

// Handler processes interception requests.
//
// Extracting the event-type from the GitHub hook event header.
//
// If we don't have a handler for this event-type, the handler returns
// the body, and a successful response, allowing unknown events through.
//
// Otherwise, they're passed to a handler to decide whether or not to
// allow the interception to complete.
func Handler(w http.ResponseWriter, r *http.Request) {
	eventType := r.Header.Get(gitHubEventHeader)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read the request body: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	h, ok := eventHandlerMap[eventType]

	if !ok {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Write(body)
		return
	}

	newBody, err := h(r, body)
	if err != nil {
		msg := fmt.Sprintf("failed handling the event: %s", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if len(newBody) > 0 {
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Write(newBody)
		return
	}

	http.Error(w, "failed interception", http.StatusPreconditionFailed)
}
