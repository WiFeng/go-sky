package elasticsearch

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"
)

var (
	testService = "testService"
)

func TestMain(m *testing.M) {

	esConf := []config.Elasticsearch{
		{
			Name: "es1",
		},
	}

	logConf := config.Log{
		Level: "info",
	}
	if _, err := log.Init(context.Background(), testService, logConf); err != nil {
		fmt.Println("Error:", err)
	}

	Init(context.Background(), testService, esConf)

	os.Exit(m.Run())
}

func TestPing(t *testing.T) {

}
