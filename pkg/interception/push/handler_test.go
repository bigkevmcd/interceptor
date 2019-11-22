package push

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-github/v28/github"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func TestHandleWithSuccess(t *testing.T) {
	event := &github.PushEvent{
		Ref: stringPtr("refs/heads/master"),
		Repo: &github.PushEventRepository{
			FullName: stringPtr("testing/testing"),
		},
	}
	r := makeRequest(t, event, "push", "master")
	body := mustMarshal(t, event)
	newBody, err := Handler(r, body)

	if err != nil {
		t.Fatal(err)
	}

	interceptedRef := gjson.GetBytes(newBody, "intercepted.ref")
	if interceptedRef.Value() != "master" {
		t.Errorf("intercepted.ref got %s, wanted %s", interceptedRef, "master")
	}
	// Delete the addition to simplify the return comparison.
	returnBody, err := sjson.DeleteBytes(newBody, "intercepted")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(returnBody, body) {
		t.Fatalf("handler got incorrect body: got %s, wanted %s", newBody, body)
	}

}

func mustMarshal(t *testing.T, e interface{}) []byte {
	body, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	return body
}
