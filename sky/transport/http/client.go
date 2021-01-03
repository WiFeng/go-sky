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
)

var (
	clientMap    = map[string]*http.Client{}
	clientConfig = map[string]config.Client{}
)

var (
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("client config is not found")
)

// Client ...
type Client struct {
	*kithttp.Client
}

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

		client := &http.Client{
			Transport: tr,
			Timeout:   cf.Timeout * unit,
		}
		clientMap[cf.Name] = client
	}
}

// NewClient ...
func NewClient(
	name string,
	method string,
	uri string,
	enc kithttp.EncodeRequestFunc,
	dec kithttp.DecodeResponseFunc,
	opt ...kithttp.ClientOption) (*Client, error) {

	cl, ok := clientMap[name]
	if !ok {
		return nil, ErrConfigNotFound
	}

	clf, ok := clientConfig[name]
	if !ok {
		return nil, ErrConfigNotFound
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

	tgt := &url.URL{
		Scheme: clf.Protocol,
		Host:   host,
		Path:   uri,
	}

	kc := kithttp.NewClient(method, tgt, enc, dec, options...)
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
