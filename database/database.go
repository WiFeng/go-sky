package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/WiFeng/go-sky/config"
	skysql "github.com/WiFeng/go-sky/database/sql"
	"github.com/WiFeng/go-sky/log"
)

var (
	dbMap    = map[string]*sql.DB{}
	dbConfig = map[string]config.Database{}
)

var (
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("database config is not found")
)

// Init ...
func Init(ctx context.Context, serviceName string, cfs []config.Database) {

	for _, cf := range cfs {
		dbConfig[cf.Name] = cf

		var db *sql.DB
		var err error
		{
			if cf.Driver == "" {
				cf.Driver = "mysql"
			}
			if cf.Charset == "" {
				cf.Charset = "utf8mb4"
			}
			if cf.InterpolateParams == "" {
				cf.InterpolateParams = "true"
			}
			if cf.DataSource == "" {
				cf.DataSource = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&interpolateParams=%s", cf.User, cf.Pass, cf.Host, cf.Port, cf.DB, cf.Charset, cf.InterpolateParams)
			}

			driverName := skysql.Register(cf.Driver)
			if db, err = sql.Open(driverName, cf.DataSource); err != nil {
				log.Fatalw(ctx, "database open error", "conf", cf, "err", err)
				continue
			}
			if err = db.PingContext(ctx); err != nil {
				log.Fatalw(ctx, "database ping error", "conf", cf, "err", err)
				continue
			}
		}

		// Dont show passwd in log
		cf.Pass = "dont show me!"
		log.Infof(ctx, "Init database [%s] %+v", cf.Name, cf)
		dbMap[cf.Name] = db
	}
}

// GetInstance ...
func GetInstance(ctx context.Context, instanceName string) (*sql.DB, error) {
	db, ok := dbMap[instanceName]
	if !ok {
		err := ErrConfigNotFound
		log.Errorw(ctx, "database.GetInstance, instanceName is not in dbMap map", "instance_name", instanceName, "err", err)
		return nil, err
	}
	return db, nil
}
