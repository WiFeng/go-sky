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

var (
	globalConfig      config.Config
	globalConfigDir   string
	globalEnvironment string
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

func init() {
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

		globalConfigDir = *configDir
		globalEnvironment = *environment
	}

	// Initialize global config
	var confFile string
	{
		if confFile, err = config.Init(globalConfigDir, globalEnvironment, &globalConfig); err != nil {
			fmt.Printf("Init config file error. path:%s, err:%v\n", confFile, err)
			os.Exit(1)
		}
	}

	// Initialzie logger
	{
		var logger log.Logger
		if logger, err = log.Init(globalConfig.Server.Log); err != nil {
			fmt.Println("Init logger error. ", err)
			os.Exit(1)
		}
		_ = logger

		// TODO:
		// defer logger.Sync()
	}

	// Initialize global tracer
	{
		var tracerCloser io.Closer
		serivceName := globalConfig.Server.Name
		if _, tracerCloser, err = trace.Init(serivceName); err != nil {
			fmt.Println("Init trace error. ", err)
			os.Exit(1)
		}
		defer tracerCloser.Close()
	}

	// Initialize pprof
	{
		pprofHost := globalConfig.Server.PProf.Host
		pporfPort := globalConfig.Server.PProf.Port
		pprof.Init(context.Background(), pprofHost, pporfPort)
	}

	// Initialize client
	{
		skyhttp.InitClient(context.Background(), globalConfig.Client)
	}
}

// LoadConfig ...
func LoadConfig(name string, conf interface{}) (err error) {
	var confFile string
	if confFile, err = config.LoadConfig(globalConfigDir, name, globalEnvironment, conf); err != nil {
		log.Errorf(context.Background(), "Load config error. path:%s, err:%v", confFile, err)
		return
	}

	log.Infof(context.Background(), "Load config successfully. path:%s", confFile)
	return
}

// LoadAppConfig ...
func LoadAppConfig(conf interface{}) error {
	return LoadConfig("app", conf)
}

// Run ...
func Run(httpHandler http.Handler) {

	httpConfig := globalConfig.Server.HTTP
	skyhttp.ListenAndServe(context.Background(), httpConfig, httpHandler)

}
