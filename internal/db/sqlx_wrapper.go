package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

const TracerName = "db"
const InstrumentVersion = "v1.0.0"
const (
	DBStatement = attribute.Key("db_statement")
	Args        = attribute.Key("db_args")
)

func recordError(span trace.Span, err error) {
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

func makeParamAttr(args []any) attribute.KeyValue {
	ss := make([]string, len(args))
	for i := range args {
		t := reflect.TypeOf(args[i])
		ss[i] = fmt.Sprintf("%s: %v", t, args[i])
	}

	return Args.StringSlice(ss)
}

func NewSqlxWrapper(db queryExecutor) *SqlxWrapper {
	tp := otel.GetTracerProvider()
	return &SqlxWrapper{
		queryExecutor: db,
		tracer:        tp.Tracer(TracerName, trace.WithInstrumentationVersion(InstrumentVersion)),
		attrs:         []attribute.KeyValue{semconv.DBSystemMySQL},
	}
}

type SqlxWrapper struct {
	queryExecutor queryExecutor
	tracer        trace.Tracer
	attrs         []attribute.KeyValue
}

func (s *SqlxWrapper) QueryxContext(ctx context.Context, query string, arg ...interface{}) (*sqlx.Rows, *sqlx.Stmt, error) {
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
	}

	spanName := "sqlx: Prepared Statement"
	ctx, span := s.tracer.Start(ctx, spanName, opts...)
	defer span.End()

	stmt, err := s.queryExecutor.PreparexContext(ctx, query)
	if err != nil {
		recordError(span, err)
		return nil, nil, err
	}

	spanName = "sqlx: Query Multiple Rows"
	ctx, spanQueryx := s.tracer.Start(ctx, spanName, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
		trace.WithAttributes(makeParamAttr(arg)),
	}...)
	defer spanQueryx.End()

	res, err := stmt.QueryxContext(ctx, arg...)
	if err != nil {
		recordError(spanQueryx, err)
		return nil, nil, err
	}

	return res, stmt, nil
}

func (s *SqlxWrapper) ExecContext(ctx context.Context, query string, arg ...interface{}) (sql.Result, error) {
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
	}

	spanName := "sqlx: Prepared Statement"
	ctx, span := s.tracer.Start(ctx, spanName, opts...)
	defer span.End()

	stmt, err := s.queryExecutor.PreparexContext(ctx, query)
	if err != nil {
		recordError(span, err)
		return nil, err
	}
	defer func(stmt *sqlx.Stmt) {
		if errCloseStmt := stmt.Close(); errCloseStmt != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "Close Prepared Statement Successfully")
		}
	}(stmt)

	spanName = "sqlx: Exec query"
	ctx, spanExec := s.tracer.Start(ctx, spanName, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
		trace.WithAttributes(makeParamAttr(arg)),
	}...)
	defer spanExec.End()

	res, err := stmt.ExecContext(ctx, arg...)
	if err != nil {
		recordError(spanExec, err)
		return nil, err
	}

	return res, nil
}

func (s *SqlxWrapper) QueryRowxContext(ctx context.Context, query string, arg ...interface{}) (*sqlx.Row, *sqlx.Stmt, error) {
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
	}

	spanName := "sqlx: Prepared Statement"
	ctx, span := s.tracer.Start(ctx, spanName, opts...)
	defer span.End()

	stmt, err := s.queryExecutor.PreparexContext(ctx, query)
	if err != nil {
		recordError(span, err)
		return nil, nil, err
	}

	spanName = "sqlx: Query Single Row"
	ctx, spanQueryx := s.tracer.Start(ctx, spanName, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
		trace.WithAttributes(makeParamAttr(arg)),
	}...)
	defer spanQueryx.End()

	res := stmt.QueryRowxContext(ctx, arg...)

	return res, stmt, nil
}
