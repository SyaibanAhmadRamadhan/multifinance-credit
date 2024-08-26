package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/repository/datastore"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/pagination"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/util/tracer"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

type callbackRows func(rows *sqlx.Rows) (err error)

type QueryRowScanType uint8

const (
	QueryRowScanTypeDefault QueryRowScanType = iota + 1
	QueryRowScanTypeStruct
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

func NewRdbms(db queryExecutor) *rdbms {
	tp := otel.GetTracerProvider()
	return &rdbms{
		db:     db,
		tracer: tp.Tracer(TracerName, trace.WithInstrumentationVersion(InstrumentVersion)),
		attrs:  []attribute.KeyValue{semconv.DBSystemMySQL},
	}
}

type rdbms struct {
	db     queryExecutor
	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

func (s *rdbms) QuerySq(ctx context.Context, query squirrel.Sqlizer, callback callbackRows) error {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return tracer.Error(err)
	}

	ctx, spanQueryx := s.tracer.Start(ctx, rawQuery, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(rawQuery)),
		trace.WithAttributes(makeParamAttr(args)),
	}...)
	defer spanQueryx.End()

	res, err := s.db.QueryxContext(ctx, rawQuery, args...)
	if err != nil {
		recordError(spanQueryx, err)
		return err
	}
	defer func() {
		if errClose := res.Close(); errClose != nil {
			recordError(spanQueryx, errClose)
		} else {
			spanQueryx.SetAttributes(attribute.String("db_closed_rows", "successfully"))
		}
	}()

	return callback(res)
}

func (s *rdbms) QuerySqPagination(ctx context.Context, countQuery, query squirrel.SelectBuilder, paginationInput pagination.PaginationInput, callback callbackRows) (
	pagination.PaginationOutput, error) {

	offset := pagination.GetOffsetValue(paginationInput.Page, paginationInput.PageSize)
	query = query.Limit(uint64(paginationInput.PageSize))
	query = query.Offset(uint64(offset))

	totalData := int64(0)
	err := s.QueryRowSq(ctx, countQuery, QueryRowScanTypeDefault, &totalData)
	if err != nil {
		return pagination.PaginationOutput{}, tracer.Error(err)
	}

	err = s.QuerySq(ctx, query, callback)
	if err != nil {
		return pagination.PaginationOutput{}, tracer.Error(err)
	}

	return pagination.CreatePaginationOutput(paginationInput, totalData), nil
}

func (s *rdbms) ExecSq(ctx context.Context, query squirrel.Sqlizer) (sql.Result, error) {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, tracer.Error(err)
	}

	ctx, spanExec := s.tracer.Start(ctx, rawQuery, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(rawQuery)),
		trace.WithAttributes(makeParamAttr(args)),
	}...)
	defer spanExec.End()

	res, err := s.db.ExecContext(ctx, rawQuery, args...)
	if err != nil {
		recordError(spanExec, err)
		return nil, err
	}

	return res, nil
}

func (s *rdbms) QueryRowSq(ctx context.Context, query squirrel.Sqlizer, scanType QueryRowScanType, dest interface{}) error {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return tracer.Error(err)
	}

	ctx, spanQueryx := s.tracer.Start(ctx, rawQuery, []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(s.attrs...),
		trace.WithAttributes(DBStatement.String(rawQuery)),
		trace.WithAttributes(makeParamAttr(args)),
	}...)
	defer spanQueryx.End()

	res := s.db.QueryRowxContext(ctx, rawQuery, args...)

	switch scanType {
	case QueryRowScanTypeStruct:
		err = res.StructScan(dest)
	default:
		err = res.Scan(dest)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.Join(err, datastore.ErrRecordNotFound)
		} else {
			recordError(spanQueryx, err)
		}

		return tracer.Error(err)
	}
	return nil
}
