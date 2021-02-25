package http

import (
	"context"
	"net/http"

	_ "net/http/pprof" // pprof

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
)

// InitPProf ...
func InitPProf(ctx context.Context, serviceName string, cfg config.PProf) {
	if cfg.Addr == "" {
		return
	}

	go func() {
		log.Infof(ctx, "Start HTTP PProf. http://%s", cfg.Addr)
		log.Fatal(ctx, http.ListenAndServe(cfg.Addr, nil))
	}()

}
