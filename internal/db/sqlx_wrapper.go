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

const TracerName = "sqlx_tracer_otel"
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

func NewRdbms(db Rdbms) *rdbms {
	tp := otel.GetTracerProvider()
	return &rdbms{
		db:     db,
		tracer: tp.Tracer(TracerName, trace.WithInstrumentationVersion(InstrumentVersion)),
		attrs:  []attribute.KeyValue{semconv.DBSystemMySQL},
	}
}

type rdbms struct {
	db     Rdbms
	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

func (s *rdbms) QueryxContext(ctx context.Context, query string, arg ...interface{}) (*sqlx.Rows, error) {
	spanName := "sqlx: Query Multiple Rows"
	ctx, spanQueryx := s.tracer.Start(ctx, spanName, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
		trace.WithAttributes(makeParamAttr(arg)),
	}...)
	defer spanQueryx.End()

	res, err := s.db.QueryxContext(ctx, query, arg...)
	if err != nil {
		recordError(spanQueryx, err)
		return nil, err
	}

	return res, nil
}

func (s *rdbms) ExecContext(ctx context.Context, query string, arg ...interface{}) (sql.Result, error) {
	spanName := "sqlx: Exec query"
	ctx, spanExec := s.tracer.Start(ctx, spanName, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
		trace.WithAttributes(makeParamAttr(arg)),
	}...)
	defer spanExec.End()

	res, err := s.db.ExecContext(ctx, query, arg...)
	if err != nil {
		recordError(spanExec, err)
		return nil, err
	}

	return res, nil
}

func (s *rdbms) QueryRowxContext(ctx context.Context, query string, arg ...interface{}) *sqlx.Row {
	spanName := "sqlx: Query Single Row"
	ctx, spanQueryx := s.tracer.Start(ctx, spanName, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(query)),
		trace.WithAttributes(makeParamAttr(arg)),
	}...)
	defer spanQueryx.End()

	res := s.db.QueryRowxContext(ctx, query, arg...)

	return res
}
