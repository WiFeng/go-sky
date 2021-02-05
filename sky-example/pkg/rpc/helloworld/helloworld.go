package helloworld

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	serviceName = "sky.helloworld"
)

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}
