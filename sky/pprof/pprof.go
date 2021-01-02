package pprof

import (
	"net/http"

	"github.com/WiFeng/go-sky/sky/log"
)

// Init ...
func Init(host string, port int) {
	if port < 1 {
		return nil
	}

	go func() {
		addr := sprinrf("%s:%d", host, port)
		log.Infof("Start HTTP PPorf. http://%s", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}()
}
