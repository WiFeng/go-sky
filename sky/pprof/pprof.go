package pprof

import (
	"context"
	"fmt"
	"net/http"

	"github.com/WiFeng/go-sky/sky/log"
)

// Init ...
func Init(ctx context.Context, host string, port int) {
	if port < 1 {
		return
	}

	go func() {
		addr := fmt.Sprintf("%s:%d", host, port)
		log.Infof(ctx, "Start HTTP PPorf. http://%s", addr)
		log.Fatal(ctx, http.ListenAndServe(addr, nil))
	}()

}
