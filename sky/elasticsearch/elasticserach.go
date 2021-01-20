package elasticsearch

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	"github.com/elastic/go-elasticsearch/v7"
)

var (
	esMap    = map[string]*elasticsearch.Client{}
	esConfig = map[string]config.Elasticsearch{}
)

var (
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("elasticsearch config is not found")
)

// Init ...
func Init(ctx context.Context, cfs []config.Elasticsearch) {

	for _, cf := range cfs {
		esConfig[cf.Name] = cf

		var cl *elasticsearch.Client
		var err error
		{
			tr := http.DefaultTransport
			if cf.CustomTranport {
				tr = &http.Transport{
					MaxIdleConnsPerHost:   cf.Transport.MaxIdleConnsPerHost,
					ResponseHeaderTimeout: cf.Transport.ResponseHeaderTimeout * time.Second,
					DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
					TLSClientConfig: &tls.Config{
						MinVersion: tls.VersionTLS11,
					},
				}
			}
			esCfg := elasticsearch.Config{
				Addresses: cf.Addresses,
				Username:  cf.Username,
				Password:  cf.Password,
				Transport: tr,
			}

			if cl, err = elasticsearch.NewClient(esCfg); err != nil {
				log.Fatalw(ctx, "elasticsearch.NewClient error", "conf", cf, "err", err)
				continue
			}
		}

		log.Infof(ctx, "Init elasticsearch [%s] %+v", cf.Name, cf)
		esMap[cf.Name] = cl
	}
}

// GetInstance ...
func GetInstance(ctx context.Context, instanceName string) (*elasticsearch.Client, error) {
	es, ok := esMap[instanceName]
	if !ok {
		err := ErrConfigNotFound
		log.Errorw(ctx, "redis.GetInstance, instanceName is not in esMap map", "instance_name", instanceName, "err", err)
		return nil, err
	}
	return es, nil
}
