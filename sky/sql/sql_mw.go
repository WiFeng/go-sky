package sql

import (
	"context"
	"database/sql"
)

var (
	// ErrDeprecatedMethod ...
	ErrDeprecatedMethod = errors.new("sql deprecated method")
)

// Conn ...
type Conn struct {
	*sql.Conn
}

// DB ...
type DB struct {
	*sql.DB
}

// Begin ...
func (db *DB) Begin() (*Tx, error) {
	return nil, ErrDeprecatedMethod
}

// BeginTx ...
func (db *DB) BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error) {

}

// Close ...
func (db *DB) Close() error {
	return db.DB.Close()
}

// Conn ...
func (db *DB) Conn(ctx context.Context) (*Conn, error) {
	return db.DB.Conn()
}

// Exec ...
func (db *DB) Exec(query string, args ...interface{}) (Result, error) {
	return nil, ErrDeprecatedMethod
}

// ExecContext ...
func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return db.DB.ExecContext(ctx, query, args...)
}

// Ping ....
func (db *DB) Ping() error {
	return ErrDeprecatedMethod
}

// PingContext ...
func (db *DB) PingContext(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

// Prepare ...
func (db *DB) Prepare(query string) (*Stmt, error) {
	return nil, ErrDeprecatedMethod
}

// PrepareContext ...
func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	return db.DB.PrepareContext(ctx, query)
}

// Query ...
func (db *DB) Query(query string, args ...interface{}) (*Rows, error) {
	return nil, ErrDeprecatedMethod
}

// QueryContext ...
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

// QueryRow ...
func (db *DB) QueryRow(query string, args ...interface{}) *Row {
	return nil, ErrDeprecatedMethod
}

// QueryRowContext ...
func (db *DDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row {
	return db.DB.QueryRowContext(ctx, query, args...)
}

// Stmt ...
type Stmt struct {
	*sql.Stmt
}

// Tx ...
type Tx struct {
	*sql.Tx
}
