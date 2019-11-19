package pullrequest

import (
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/go-github/v28/github"
)

func TestHandleWithSuccess(t *testing.T) {
	event := &github.PullRequestEvent{
		Action: stringPtr("open"),
		Repo: &github.Repository{
			FullName: stringPtr("testing/testing"),
		},
	}
	r, body := makeRequest(t, event, "pull_request", "open")
	w := httptest.NewRecorder()

	InterceptionHandler(w, r)

	resp := w.Result()
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type incorrect, got %s, wanted %s", ct, "application/json")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(respBody, body) {
		t.Fatalf("decoded response: got %+v, wanted %+v\n", respBody, body)
	}
}
