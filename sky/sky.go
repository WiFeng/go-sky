package sky

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-lib/metrics/prometheus"

	skyhttp "github.com/WiFeng/go-sky/sky/transport/http"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

var (
	globalConfig      config.Config
	globalConfigDir   string
	globalConfigFile  string
	globalEnvironment string
)

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
		globalConfigFile = confFile
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
		if _, tracerCloser, err = initTrace(serivceName); err != nil {
			fmt.Println("Init trace error. ", err)
			os.Exit(1)
		}

		// TODO:
		_ = tracerCloser
		// defer tracerCloser.Close()
	}

	// Initialize pprof
	{
		pprofHost := globalConfig.Server.PProf.Host
		pporfPort := globalConfig.Server.PProf.Port
		initPProf(context.Background(), pprofHost, pporfPort)
	}

	// Initialize client
	{
		skyhttp.InitClient(context.Background(), globalConfig.Client)
	}

	log.Infow(context.Background(), "Load config successfully", "path", globalConfigFile, "env", globalEnvironment)
}

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

func initPProf(ctx context.Context, host string, port int) {
	if port < 1 {
		return
	}

	go func() {
		addr := fmt.Sprintf("%s:%d", host, port)
		log.Infof(ctx, "Start HTTP PPorf. http://%s", addr)
		log.Fatal(ctx, http.ListenAndServe(addr, nil))
	}()

}

func initTrace(serviceName string) (opentracing.Tracer, io.Closer, error) {
	metricsFactory := prometheus.New()

	//logger := log.GetDefaultLogger()
	loggerOption := jaegerconfig.Logger(jaegerlog.DebugLogAdapter(jaeger.StdLogger))
	tracer, tracerCloser, err := jaegerconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			// LogSpans:          true,
			// CollectorEndpoint: "http://localhost:14268/api/traces",
			LocalAgentHostPort:  "localhost:6831",
			BufferFlushInterval: time.Second * 1,
		},
	}.NewTracer(
		jaegerconfig.Metrics(metricsFactory),
		loggerOption,
	)
	opentracing.InitGlobalTracer(tracer)
	return tracer, tracerCloser, err
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

// LoadConfig ...
func LoadConfig(name string, conf interface{}) (err error) {
	var confFile string
	if confFile, err = config.LoadConfig(globalConfigDir, name, globalEnvironment, conf); err != nil {
		log.Errorw(context.Background(), "Load config error", "path", confFile, "err", err)
		return
	}

	log.Infow(context.Background(), "Load config successfully", "path", confFile)
	return
}

// LoadAppConfig ...
func LoadAppConfig(conf interface{}) error {
	return LoadConfig("app", conf)
}

// RunHTTPServer ...
func RunHTTPServer(handler http.Handler) {
	httpConfig := globalConfig.Server.HTTP
	skyhttp.ListenAndServe(context.Background(), httpConfig, handler)
}

// Run ...
func Run(handler http.Handler) {
	RunHTTPServer(handler)
}
