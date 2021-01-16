package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	"github.com/WiFeng/go-sky/sky/middleware"
)

var (
	clientMap    = map[string]kithttp.HTTPClient{}
	clientConfig = map[string]config.Client{}
)

var (
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("client config is not found")
)

// InitClient ...
func InitClient(ctx context.Context, cfs []config.Client) {
	for _, cf := range cfs {
		clientConfig[cf.Name] = cf

		var tr http.RoundTripper

		var unit = time.Second
		if cf.MillSecUnit {
			unit = time.Millisecond
		}

		if cf.CustomTranport {
			tr = &http.Transport{

				MaxConnsPerHost:     cf.MaxConnsPerHost,
				MaxIdleConns:        cf.MaxIdleConns,
				MaxIdleConnsPerHost: cf.MaxIdleConnsPerHost,

				IdleConnTimeout:       cf.IdleConnTimeout * unit,
				TLSHandshakeTimeout:   cf.TLSHandshakeTimeout * unit,
				ExpectContinueTimeout: cf.ExpectContinueTimeout * unit,
				ResponseHeaderTimeout: cf.ResponseHeaderTimeout * unit,

				DisableKeepAlives:  cf.DisableKeepAlives,
				DisableCompression: cf.DisableCompression,
			}
		} else {
			tr = http.DefaultTransport
		}

		cl := &middleware.Client{
			Client: &http.Client{
				Transport: tr,
				Timeout:   cf.Timeout * unit,
			},
		}
		cl.Use(middleware.HTTPClientLoggingMiddleware)
		cl.Use(middleware.HTTPClientTracingMiddleware)
		clientMap[cf.Name] = cl
	}
}

// Client ...
type Client struct {
	*kithttp.Client
}

// NewClient ...
func NewClient(
	ctx context.Context,
	serviceName string,
	method string,
	uri string,
	enc kithttp.EncodeRequestFunc,
	dec kithttp.DecodeResponseFunc,
	opt ...kithttp.ClientOption) (*Client, error) {

	cl, ok := clientMap[serviceName]
	if !ok {
		err := ErrConfigNotFound
		log.Errorw(ctx, "http.NewClient, serviceName is not in clientMap map",
			"service_name", serviceName, "method", method, "uri", uri, "err", err)
		return nil, err
	}

	clf, ok := clientConfig[serviceName]
	if !ok {
		err := ErrConfigNotFound
		log.Errorw(ctx, "http.NewClient, serviceName is not in clientConfig map",
			"service_name", serviceName, "method", method, "uri", uri, "err", err)
		return nil, err
	}

	options := []kithttp.ClientOption{
		kithttp.SetClient(cl),
	}

	if opt != nil {
		options = append(options, opt...)
	}

	host := clf.Host
	if clf.Port > 0 {
		host = fmt.Sprintf("%s:%d", clf.Host, clf.Port)
	}

	targetURL, err := url.Parse(uri)
	if err != nil {
		log.Errorw(ctx, "http.NewClient, url.Parse error",
			"service_name", serviceName, "method", method, "uri", uri, "err", err)
		return nil, err
	}

	targetURL.Scheme = clf.Protocol
	targetURL.Host = host

	kc := kithttp.NewClient(method, targetURL, enc, dec, options...)
	c := &Client{
		kc,
	}

	return c, nil
}

// Endpoint ...
func (c Client) Endpoint() kitendpoint.Endpoint {
	e := c.Client.Endpoint()
	// e = opentracing.TraceClient(otTracer, "Sum")(e)

	return e
}
