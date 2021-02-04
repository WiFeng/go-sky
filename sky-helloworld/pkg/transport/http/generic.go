package http

import (
	"context"
	"encoding/json"
	"net/http"
)

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
