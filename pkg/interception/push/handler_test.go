package push

import (
	"reflect"
	"testing"

	"github.com/google/go-github/v28/github"
)

func TestHandleWithSuccess(t *testing.T) {
	event := &github.PushEvent{
		Ref: stringPtr("refs/heads/master"),
		Repo: &github.PushEventRepository{
			FullName: stringPtr("testing/testing"),
		},
	}
	r, body := makeRequest(t, event, "push", "master")

	newBody, err := Handler(r, body)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(newBody, body) {
		t.Fatalf("handler got incorrect body: got %s, wanted %s", newBody, body)
	}

}