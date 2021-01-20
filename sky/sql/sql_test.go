package sql

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
	dbConf := []config.DB{
		{
			Name: "db1",
			Host: "127.0.0.1",
			DB:   "test",
			Port: 3306,
			User: "root",
			Pass: "",
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

	Init(context.Background(), dbConf)

	os.Exit(m.Run())
}

func TestSelect(t *testing.T) {
	var ctx = context.Background()
	db, err := GetInstance(ctx, "db11")
	if err != nil {
		t.Error(err)
		return
	}

	var got string
	err = db.QueryRowContext(ctx, "SELECT OK").Scan(&got)
	if err != nil {
		t.Error(err)
		return
	}
	if got != "OK" {
		t.Errorf("db.QueryRowContext = %s; want OK", got)
		return
	}
}
