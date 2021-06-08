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

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
)

const (
	clientContext clientContextKey = 0
)

var (
	clientMap    = map[string]*http.Client{}
	clientConfig = map[string]config.Client{}
)

var (
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("client config is not found")
)

type clientContextKey int
type clientContextVal struct {
	peerName string
}

// InitClient ...
func InitClient(ctx context.Context, peerName string, cfs []config.Client) {
	for _, cf := range cfs {
		clientConfig[cf.Name] = cf

		tr := NewRoundTripperFromConfig(cf.Transport)
		tr.Use(RoundTripperMetricsMiddleware)
		tr.Use(RoundTripperTracingMiddleware)
		tr.Use(RoundTripperLoggingMiddleware)

		var timeout time.Duration
		if cf.Timeout > 0 {
			timeout = cf.Timeout * time.Second
		}
		if cf.TimeoutMillSec > 0 {
			timeout = cf.TimeoutMillSec * time.Millisecond
		}
		cl := &http.Client{
			Transport: tr,
			Timeout:   timeout,
		}
		clientMap[cf.Name] = cl
	}
}

// Client ...
type Client struct {
	*kithttp.Client
	peerName string
}

// NewClient ...
func NewClient(
	ctx context.Context,
	peerName string,
	method string,
	uri string,
	enc kithttp.EncodeRequestFunc,
	dec kithttp.DecodeResponseFunc,
	opt ...kithttp.ClientOption) (*Client, error) {

	cl, ok := clientMap[peerName]
	if !ok {
		err := ErrConfigNotFound
		log.Errorw(ctx, "http.NewClient, peerName is not in clientMap map",
			"service_name", peerName, "method", method, "uri", uri, "err", err)
		return nil, err
	}

	clf, ok := clientConfig[peerName]
	if !ok {
		err := ErrConfigNotFound
		log.Errorw(ctx, "http.NewClient, peerName is not in clientConfig map",
			"service_name", peerName, "method", method, "uri", uri, "err", err)
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
			"service_name", peerName, "method", method, "uri", uri, "err", err)
		return nil, err
	}

	targetURL.Scheme = clf.Protocol
	targetURL.Host = host

	kc := kithttp.NewClient(method, targetURL, enc, dec, options...)
	c := &Client{
		kc,
		peerName,
	}

	return c, nil
}

// Endpoint ...
func (c Client) Endpoint() kitendpoint.Endpoint {
	// e := c.Client.Endpoint()
	// return e
	e := c.Client.Endpoint()
	f := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ctx = context.WithValue(ctx, clientContext, clientContextVal{c.peerName})
		return e(ctx, request)
	}
	return f
}
