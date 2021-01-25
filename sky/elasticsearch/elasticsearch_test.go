package elasticsearch

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {

	esConf := []config.Elasticsearch{
		{
			Name: "es1",
		},
	}

	logConf := config.Log{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
	}
	if _, err := log.Init(logConf); err != nil {
		fmt.Println("Error:", err)
	}

	Init(context.Background(), esConf)

	os.Exit(m.Run())
}

func TestPing(t *testing.T) {

}
