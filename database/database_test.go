package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/WiFeng/go-sky/config"
	"github.com/WiFeng/go-sky/log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	testName    = "db1"
	testService = "testService"
)

func TestMain(m *testing.M) {
	dbConf := []config.Database{
		{
			Name: testName,
			Host: "127.0.0.1",
			DB:   "test",
			Port: 3306,
			User: "root",
			Pass: "123456",
		},
	}

	logConf := config.Log{
		Level: "info",
	}
	if _, err := log.Init(context.Background(), testService, logConf); err != nil {
		fmt.Println("Error:", err)
	}

	Init(context.Background(), testService, dbConf)

	os.Exit(m.Run())
}

func TestSelect(t *testing.T) {
	var ctx = context.Background()
	db, err := GetInstance(ctx, testName)
	if err != nil {
		t.Error(err)
		return
	}

	var got string
	err = db.QueryRowContext(ctx, "SELECT 'OK'").Scan(&got)
	if err != nil {
		t.Error(err)
		return
	}
	if got != "OK" {
		t.Errorf("db.QueryRowContext = %s; want OK", got)
		return
	}
}

func TestStmtQueryContext(t *testing.T) {
	var ctx = context.Background()
	db, err := GetInstance(ctx, testName)
	if err != nil {
		t.Error(err)
		return
	}

	var got string
	stmt, err := db.PrepareContext(ctx, "SELECT ?")
	if err != nil {
		t.Error(err)
		return
	}

	stmt.QueryRowContext(ctx, "OK").Scan(&got)
	if err != nil {
		t.Error(err)
		return
	}

	if got != "OK" {
		t.Errorf("db.QueryRowContext = %s; want OK", got)
		return
	}
}

func TestStmtInsertContext(t *testing.T) {
	var ctx = context.Background()
	db, err := GetInstance(ctx, testName)
	if err != nil {
		t.Error(err)
		return
	}

	stmt, err := db.PrepareContext(ctx, "INSERT INTO te(name, salary) values(?,?)")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = stmt.ExecContext(ctx, "Lucy", 1200)
	if err != nil {
		t.Error(err)
		return
	}
}
