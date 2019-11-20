package interception

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandlerWithUnknownEventTypeReturnsBody(t *testing.T) {
	_, ok := eventHandlerMap["unknown"]
	if ok {
		t.Fatal("unknown event has a handler")
	}
	testBody := []byte(`{}`)
	r, _ := http.NewRequest("POST", "/", bytes.NewReader(testBody))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(gitHubEventHeader, "unknown")
	w := httptest.NewRecorder()

	Handler(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code, got %d, wanted %d", resp.StatusCode, http.StatusOK)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type incorrect, got %s, wanted %s", ct, "application/json")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(respBody, testBody) {
		t.Errorf("decoded response: got %+v, wanted %+v\n", respBody, testBody)
	}

}

func TestSuccessfulResponse(t *testing.T) {
	testResponse := []byte(`testing`)
	eventHandlerMap["pull_request"] = func(r *http.Request, body []byte) ([]byte, error) {
		return testResponse, nil
	}
	r := makePullRequestRequest(t, []byte(`{}`))
	w := httptest.NewRecorder()

	Handler(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code, got %d, wanted %d", resp.StatusCode, http.StatusOK)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type incorrect, got %s, wanted %s", ct, "application/json")
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(respBody, testResponse) {
		t.Errorf("decoded response: got %+v, wanted %+v\n", respBody, testResponse)
	}
}

func TestHandlerNoResponse(t *testing.T) {
	eventHandlerMap["pull_request"] = func(r *http.Request, body []byte) ([]byte, error) {
		return nil, nil
	}
	r := makePullRequestRequest(t, []byte(`{}`))
	w := httptest.NewRecorder()

	Handler(w, r)

	resp := w.Result()
	if resp.StatusCode != http.StatusPreconditionFailed {
		t.Errorf("unexpected status code, got %d, wanted %d", resp.StatusCode, http.StatusPreconditionFailed)
	}
}

func TestErrorResponse(t *testing.T) {
	eventHandlerMap["pull_request"] = func(r *http.Request, body []byte) ([]byte, error) {
		return nil, errors.New("test error")
	}
	r := makePullRequestRequest(t, []byte(`{}`))
	w := httptest.NewRecorder()

	Handler(w, r)

	resp := w.Result()
	wantedStatus := http.StatusInternalServerError
	if resp.StatusCode != wantedStatus {
		t.Errorf("unexpected status code, got %d, wanted %d", resp.StatusCode, wantedStatus)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	wantedMsg := "failed handling the event: test error\n"
	if string(body) != wantedMsg {
		t.Fatalf("unexpected error message: got %s, wanted %s", body, wantedMsg)
	}

}

func makePullRequestRequest(t *testing.T, body []byte) *http.Request {
	r, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(gitHubEventHeader, "pull_request")
	return r
}
