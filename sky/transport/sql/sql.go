package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/WiFeng/go-sky/sky/config"
	"github.com/WiFeng/go-sky/sky/log"

	_ "github.com/go-sql-driver/mysql" // mysql
)

var (
	dbMap    = map[string]*sql.DB{}
	dbConfig = map[string]config.DB{}
)

var (
	// ErrConfigNotFound ...
	ErrConfigNotFound = errors.New("database config is not found")
)

// Init ...
func Init(ctx context.Context, cfs []config.DB) {

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

			if db, err = sql.Open(cf.Driver, fmt.Sprintf("%s:%s@/%s?charset=%s", cf.User, cf.Pass, cf.DB, cf.Charset)); err != nil {
				log.Fatalw(ctx, "database open error", "conf", cf, "err", err)
				continue
			}
			if err = db.PingContext(ctx); err != nil {
				log.Fatalw(ctx, "database ping error", "conf", cf, "err", err)
				continue
			}
		}

		log.Infof(ctx, "Init database [%s] %+v", cf.Name, cf)
		dbMap[cf.Name] = db
	}
}

// GetInstance ...
func GetInstance(ctx context.Context, instanceName string) (*sql.DB, error) {
	db, ok := dbMap[instanceName]
	if !ok {
		err := ErrConfigNotFound
		log.Errorw(ctx, "redis.GetInstance, instanceName is not in dbMap map", "instance_name", instanceName, "err", err)
		return nil, err
	}
	return db, nil
}
