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

	var parentSpan opentracing.Span
	var childSpan opentracing.Span

	defer func() {
		if childSpan == nil {
			return
		}
		if err != nil {
			opentracingext.Error.Set(childSpan, true)
			childSpan.SetTag("db.error", err.Error())
			childSpan.Finish()
			return
		}
		childSpan.Finish()
	}()

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			"sql.Ping",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tag{Key: string(opentracingext.DBType), Value: "sql"},
			opentracing.Tag{Key: string(opentracingext.Component), Value: "database"},
			opentracingext.SpanKindRPCClient,
		)
	}

	// ============================================
	pinger, ok := c.base.(driver.Pinger)
	if !ok {
		return ErrUnsupportedMethod
	}
	err = pinger.Ping(ctx)
	return
}

func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (tx driver.Tx, err error) {

	var parentSpan opentracing.Span
	var childSpan opentracing.Span

	defer func() {
		if childSpan == nil {
			return
		}
		if err != nil {
			opentracingext.Error.Set(childSpan, true)
			childSpan.SetTag("db.error", err.Error())
			childSpan.Finish()
			return
		}
		childSpan.Finish()
	}()

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			"sql.BeginTx",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tag{Key: string(opentracingext.DBType), Value: "sql"},
			opentracing.Tag{Key: string(opentracingext.Component), Value: "database"},
			opentracingext.SpanKindRPCClient,
		)
	}

	// ============================================
	connBegin, ok := c.base.(driver.ConnBeginTx)
	if !ok {
		return nil, ErrUnsupportedMethod
	}
	tx, err = connBegin.BeginTx(ctx, opts)
	return
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
			childSpan.SetTag("db.error", err.Error())
			childSpan.Finish()
			return
		}
		childSpan.Finish()
	}()

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			"sql.QueryContext",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tag{Key: "db.query", Value: query},
			opentracing.Tag{Key: "db.args", Value: args},
			opentracing.Tag{Key: string(opentracingext.DBType), Value: "sql"},
			opentracing.Tag{Key: string(opentracingext.Component), Value: "database"},
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

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (result driver.Result, err error) {

	var parentSpan opentracing.Span
	var childSpan opentracing.Span

	defer func() {
		if childSpan == nil {
			return
		}
		if err != nil {
			opentracingext.Error.Set(childSpan, true)
			childSpan.SetTag("db.error", err.Error())
			childSpan.Finish()
			return
		}
		childSpan.Finish()
	}()

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			"sql.ExecContext",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tag{Key: "db.query", Value: query},
			opentracing.Tag{Key: "db.args", Value: args},
			opentracing.Tag{Key: string(opentracingext.DBType), Value: "sql"},
			opentracing.Tag{Key: string(opentracingext.Component), Value: "database"},
			opentracingext.SpanKindRPCClient,
		)
	}

	// ============================================
	execer, ok := c.base.(driver.ExecerContext)
	if !ok {
		return nil, ErrUnsupportedMethod
	}
	result, err = execer.ExecContext(ctx, query, args)
	return
}

func (c *conn) PrepareContext(ctx context.Context, query string) (s driver.Stmt, err error) {
	var parentSpan opentracing.Span
	var childSpan opentracing.Span

	defer func() {
		if childSpan == nil {
			return
		}
		if err != nil {
			opentracingext.Error.Set(childSpan, true)
			childSpan.SetTag("db.error", err.Error())
			childSpan.Finish()
			return
		}
		childSpan.Finish()
	}()

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			"sql.PrepareContext",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tag{Key: "db.query", Value: query},
			opentracing.Tag{Key: string(opentracingext.DBType), Value: "sql"},
			opentracing.Tag{Key: string(opentracingext.Component), Value: "database"},
			opentracingext.SpanKindRPCClient,
		)
	}

	// ============================================
	connPrepare, ok := c.base.(driver.ConnPrepareContext)
	if !ok {
		return nil, ErrUnsupportedMethod
	}
	basestmt, err := connPrepare.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	s = &stmt{
		base: basestmt,
	}
	return
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

	var parentSpan opentracing.Span
	var childSpan opentracing.Span

	defer func() {
		if childSpan == nil {
			return
		}
		if err != nil {
			opentracingext.Error.Set(childSpan, true)
			childSpan.SetTag("db.error", err.Error())
			childSpan.Finish()
			return
		}
		childSpan.Finish()
	}()

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			"sql.stmt.QueryContext",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tag{Key: "db.args", Value: args},
			opentracing.Tag{Key: string(opentracingext.DBType), Value: "sql"},
			opentracing.Tag{Key: string(opentracingext.Component), Value: "database"},
			opentracingext.SpanKindRPCClient,
		)
	}

	// ============================================
	queryer, ok := s.base.(driver.StmtQueryContext)
	if !ok {
		return nil, ErrUnsupportedMethod
	}
	rows, err = queryer.QueryContext(ctx, args)
	return
}

func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (result driver.Result, err error) {

	var parentSpan opentracing.Span
	var childSpan opentracing.Span

	defer func() {
		if childSpan == nil {
			return
		}
		if err != nil {
			opentracingext.Error.Set(childSpan, true)
			childSpan.SetTag("db.error", err.Error())
			childSpan.Finish()
			return
		}
		childSpan.Finish()
	}()

	if parentSpan = opentracing.SpanFromContext(ctx); parentSpan != nil {
		childSpan = parentSpan.Tracer().StartSpan(
			"sql.stmt.ExecContext",
			opentracing.ChildOf(parentSpan.Context()),
			opentracing.Tag{Key: "db.args", Value: args},
			opentracing.Tag{Key: string(opentracingext.DBType), Value: "sql"},
			opentracing.Tag{Key: string(opentracingext.Component), Value: "database"},
			opentracingext.SpanKindRPCClient,
		)
	}

	// ============================================
	execer, ok := s.base.(driver.StmtExecContext)
	if !ok {
		return nil, ErrUnsupportedMethod
	}
	result, err = execer.ExecContext(ctx, args)
	return
}
