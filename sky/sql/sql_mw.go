package sql

import (
	"context"
	"database/sql"
)

var (
	// ErrSQLDeprecatedMethod ...
	ErrSQLDeprecatedMethod = errors.new("sql deprecated method")
)

// SQLConn ...
type SQLConn struct {
	*sql.Conn
}

// SQLDB ...
type SQLDB struct {
	*sql.DB
}

// Begin ...
func (db *SQLDB) Begin() (*Tx, error) {
	return nil, ErrSQLDeprecatedMethod
}

// BeginTx ...
func (db *SQLDB) BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error) {

}

// Close ...
func (db *SQLDB) Close() error {
	return db.DB.Close()
}

// Conn ...
func (db *SQLDB) Conn(ctx context.Context) (*Conn, error) {
	return db.DB.Conn()
}

// Exec ...
func (db *SQLDB) Exec(query string, args ...interface{}) (Result, error) {
	return nil, ErrSQLDeprecatedMethod
}

// ExecContext ...
func (db *SQLDB) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return db.DB.ExecContext(ctx, query, args...)
}

// Ping ....
func (db *SQLDB) Ping() error {
	return ErrSQLDeprecatedMethod
}

// PingContext ...
func (db *SQLDB) PingContext(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

// Prepare ...
func (db *SQLDB) Prepare(query string) (*Stmt, error) {
	return nil, ErrSQLDeprecatedMethod
}

// PrepareContext ...
func (db *SQLDB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	return db.DB.PrepareContext(ctx, query)
}

// Query ...
func (db *SQLDB) Query(query string, args ...interface{}) (*Rows, error) {
	return nil, ErrSQLDeprecatedMethod
}

// QueryContext ...
func (db *SQLDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

// QueryRow ...
func (db *SQLDB) QueryRow(query string, args ...interface{}) *Row {
	return nil, ErrSQLDeprecatedMethod
}

// QueryRowContext ...
func (db *DSQLDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}

// SQLStmt ...
type SQLStmt struct {
	*sql.Stmt
}

// SQLTx ...
type SQLTx struct {
	*sql.Tx
}
