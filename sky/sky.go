package sky

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	"github.com/WiFeng/go-sky/sky/pprof"
	"github.com/WiFeng/go-sky/sky/trace"
	skyhttp "github.com/WiFeng/go-sky/sky/transport/http"
)

func initFlag() (*string, *string, error) {
	fs := flag.NewFlagSet("short-url", flag.ExitOnError)

	var (
		configDir   = fs.String("conf", "./conf/", "Config Directory")
		environment = fs.String("env", "development", "Runing environment")
	)

	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	err := fs.Parse(os.Args[1:])

	return configDir, environment, err
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

// Run ...
func Run(httpHandler http.Handler) {

	// Initialize flogs
	var err error
	var configDir *string
	var environment *string
	{
		configDir, environment, err = initFlag()
		if err != nil {
			fmt.Println("Init flag error. ", err)
			os.Exit(1)
		}
	}

	// Initialize global config
	var globalConfig config.Config
	{
		if err = config.Init(*configDir, *environment, &globalConfig); err != nil {
			fmt.Println("Init config file error. ", err)
			os.Exit(1)
		}

	}

	// Initialzie logger
	{
		var logger log.Logger
		if logger, err = log.Init(globalConfig.Server.Log); err != nil {
			fmt.Println("Init config file error. ", err)
			os.Exit(1)
		}
		defer logger.Sync()
	}

	// Initialize global tracer
	{
		var tracerCloser io.Closer
		serivceName := globalConfig.Server.Name
		if _, tracerCloser, err = trace.Init(serivceName); err != nil {
			fmt.Println("Init config file error. ", err)
			os.Exit(1)
		}
		defer tracerCloser.Close()
	}

	var ctx context.Context
	{
		ctx = context.Background()
	}

	// Initialize pprof
	{
		pprofHost := globalConfig.Server.PProf.Host
		pporfPort := globalConfig.Server.PProf.Port
		pprof.Init(ctx, pprofHost, pporfPort)
	}

	// Start listen and serve
	{
		httpConfig := globalConfig.Server.HTTP
		skyhttp.ListenAndServe(ctx, httpConfig, httpHandler)
	}

}
