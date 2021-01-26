package driver

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/opentracing/opentracing-go"
	opentracingext "github.com/opentracing/opentracing-go/ext"
)

var (
	// ErrDeprecatedMethod ...
	ErrDeprecatedMethod = errors.New("sql deprecated method")
	// ErrUnsupportedMethod ...
	ErrUnsupportedMethod = errors.New("sql unsupported method")
)

// Conn ...
type conn struct {
	base        driver.Conn
	pinger      driver.Pinger
	execer      driver.ExecerContext
	queryer     driver.QueryerContext
	connPrepare driver.ConnPrepareContext
	connBegin   driver.ConnBeginTx
}

func (c *conn) Begin() (driver.Tx, error) {
	return nil, ErrDeprecatedMethod
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return nil, ErrDeprecatedMethod
}

func (c *conn) Close() error {
	return c.base.Close()
}

func (c *conn) Ping(ctx context.Context) (err error) {
	if pinger, ok := c.base.(driver.Pinger); ok {
		return pinger.Ping(ctx)
	}
	return nil
}

func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if connBegin, ok := c.base.(driver.ConnBeginTx); ok {
		return connBegin.BeginTx(ctx, opts)
	}
	return nil, nil
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {

	var parentSpan opentracing.Span
	var childSpan opentracing.Span

	defer func() {
		if childSpan == nil {
			return
		}
		if err != nil {
			opentracingext.Error.Set(childSpan, true)
			childSpan.SetTag("http.error", err.Error())
			childSpan.Finish()
			return
		}
		childSpan.Finish()
	}()

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			"sql.QueryContext",
			opentracing.ChildOf(parentSpan.Context()),
			opentracingext.SpanKindRPCClient,
		)
	}

	// ============================================
	queryer, ok := c.base.(driver.QueryerContext)
	if !ok {
		return nil, ErrUnsupportedMethod
	}
	rows, err = queryer.QueryContext(ctx, query, args)
	return
}

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if execer, ok := c.base.(driver.ExecerContext); ok {
		return execer.ExecContext(ctx, query, args)
	}
	return nil, nil
}

func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if connPrepare, ok := c.base.(driver.ConnPrepareContext); ok {
		basestmt, err := connPrepare.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		stmt := &stmt{
			base: basestmt,
		}
		return stmt, nil
	}
	return nil, nil
}

type stmt struct {
	base    driver.Stmt
	queryer driver.StmtQueryContext
	execer  driver.StmtExecContext
}

func (s *stmt) Close() error {
	return s.base.Close()
}

func (s *stmt) NumInput() int {
	return s.base.NumInput()
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, ErrDeprecatedMethod
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, ErrDeprecatedMethod
}

func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	if queryer, ok := s.base.(driver.StmtQueryContext); ok {
		return queryer.QueryContext(ctx, args)
	}
	return
}

func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (result driver.Result, err error) {
	if execer, ok := s.base.(driver.StmtExecContext); ok {
		return execer.ExecContext(ctx, args)
	}
	return
}
