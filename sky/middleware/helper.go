package middleware

import (
	"net/http"
	"time"

	"github.com/WiFeng/go-sky/sky/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// HTTPResponseWriter ...
type HTTPResponseWriter struct {
	http.ResponseWriter
	statusCode int
	reqBody    []byte
	respBody   []byte
}

func (w *HTTPResponseWriter) Write(b []byte) (int, error) {
	w.respBody = append(w.respBody, b...)
	return w.ResponseWriter.Write(b)
}

// WriteHeader ...
func (w *HTTPResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// ClientDoFunc ...
type ClientDoFunc func(*http.Request) (*http.Response, error)

// Do ...
func (c ClientDoFunc) Do(req *http.Request) (*http.Response, error) {
	return c(req)
}

// Client ...
type Client struct {
	*http.Client
	middlewares []ClientMiddleware
}

// Use ...
func (c *Client) Use(mwf ...ClientMiddlewareFunc) {
	for _, fn := range mwf {
		c.middlewares = append(c.middlewares, fn)
	}
}

// Do ...
func (c Client) Do(req *http.Request) (*http.Response, error) {
	var cl = kithttp.HTTPClient(c.Client)
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		cl = c.middlewares[i].Middleware(cl)
	}
	return cl.Do(req)
}

// ClientMiddleware ...
type ClientMiddleware interface {
	Middleware(kithttp.HTTPClient) kithttp.HTTPClient
}

// ClientMiddlewareFunc ...
type ClientMiddlewareFunc func(kithttp.HTTPClient) kithttp.HTTPClient

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw ClientMiddlewareFunc) Middleware(httpClient kithttp.HTTPClient) kithttp.HTTPClient {
	return mw(httpClient)
}

// HTTPClientDoMiddleware ...
func HTTPClientDoMiddleware(next kithttp.HTTPClient) kithttp.HTTPClient {
	return ClientDoFunc(func(req *http.Request) (*http.Response, error) {
		ctx := req.Context()

		defer func(begin time.Time) {
			log.Infow(ctx, "", "request_time", time.Since(begin).Microseconds())
		}(time.Now())
		return next.Do(req)
	})
}
