package interception

import "net/http"

// InterceptionFunc returns the response body, and possibly an error.
// If the response body is nil, this will be returned to the client as an error,
// indicating that it should not continue.
type InterceptionFunc func(r *http.Request, body []byte) ([]byte, error)
