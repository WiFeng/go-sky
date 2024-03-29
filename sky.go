package sky

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/helper"
	"github.com/WiFeng/go-sky/log"
	"github.com/WiFeng/go-sky/trace"

	skydb "github.com/WiFeng/go-sky/database"
	skyes "github.com/WiFeng/go-sky/elasticsearch"
	skyhttp "github.com/WiFeng/go-sky/http"
	skykafka "github.com/WiFeng/go-sky/kafka"
	skymetrics "github.com/WiFeng/go-sky/metrics"
	skyredis "github.com/WiFeng/go-sky/redis"
)

var (
	globalConfigDir   string
	globalConfigFile  string
	globalEnvironment string

	globalConfig config.Config
)

func init() {

	var err error
	var ctx = context.Background()

	// Initialize flogs
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

	// Initialzie logger and trace
	log.Init(ctx, globalConfig.Server.Name, globalConfig.Server.Log)
	trace.Init(ctx, globalConfig.Server.Name, globalConfig.Server.Trace)

	// Initialzie supported components
	skydb.Init(ctx, globalConfig.Server.Name, globalConfig.Database)
	skyredis.Init(ctx, globalConfig.Server.Name, globalConfig.Redis)
	skyes.Init(ctx, globalConfig.Server.Name, globalConfig.Elasticsearch)
	skykafka.Init(ctx, globalConfig.Server.Name, globalConfig.Kafka)
	skymetrics.Init(ctx, globalConfig.Server.Name, globalConfig.Server.Metrics)
	skyhttp.InitClient(ctx, globalConfig.Server.Name, globalConfig.Client)
	skyhttp.InitPProf(ctx, globalConfig.Server.Name, globalConfig.Server.PProf)

	log.Infow(ctx, "Load config successfully", "path", globalConfigFile, "env", globalEnvironment)
}

func initFlag() (*string, *string, error) {
	fs := flag.NewFlagSet("go-sky", flag.ExitOnError)

	var (
		configDir   = fs.String("conf", "./conf/", "Config Directory")
		environment = fs.String("env", "", "Runing environment")
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

func SafelyDo(ctx context.Context, f func()) (err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Errorw(ctx, "panic error", "err", panicErr)
			err = errors.New("panic error")
		}
	}()
	f()
	return
}

// RegisterTask ...
func RegisterTask(f func(), df func(), sync bool) {
	ctx := context.Background()
	if sync {
		if SafelyDo(ctx, f) != nil {
			log.Panicf(context.Background(), "panic error")
		}
	} else {
		go SafelyDo(ctx, f)
	}

	helper.AddDeferFunc(df)
}

// RunHTTPServer ...
func RunHTTPServer(handler http.Handler) {
	// listen
	skyhttp.ListenAndServe(context.Background(), globalConfig.Server.HTTP, handler)

	// do something of the clearup
	helper.RunDeferFunc()
}

// Run ...
func Run(handler http.Handler) {
	RunHTTPServer(handler)
}
