package presentation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/multifinance-credit/internal/presentation/restapi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"strconv"
	"time"
)

var tracerName = "response-writer"
var otelTracer = otel.Tracer(tracerName)

type Option func(*ResponseWriter)

func WithLogRequestBody(log bool) Option {
	return func(e *ResponseWriter) {
		e.logReqBody = log
	}
}

func WithLogResponseBody(log bool) Option {
	return func(e *ResponseWriter) {
		e.logRespBody = log
	}
}

func WithLogParams(log bool) Option {
	return func(e *ResponseWriter) {
		e.logParams = log
	}
}

type ResponseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	logParams   bool
	logRespBody bool
	logReqBody  bool
	buffer      *bytes.Buffer
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) Write(body []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(body)
	rw.size = size
	if rw.logRespBody {
		rw.buffer = new(bytes.Buffer)
		rw.buffer.Write(body)
	}
	return size, err
}

func withOtel(next http.HandlerFunc, opts ...Option) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now().UTC()

		recorder := &ResponseWriter{
			ResponseWriter: writer,
			logParams:      true,
			logRespBody:    true,
			logReqBody:     true,
		}

		for _, opt := range opts {
			opt(recorder)
		}

		if recorder.logParams {
			request = queryParamToSpan(request, request.URL.Query())
		}

		if recorder.logReqBody && (request.Method == http.MethodPost || request.Method == http.MethodPut) {
			var err error
			request, err = addRequestBodyToSpan(request)
			if err != nil {
				restapi.Error(recorder, request, http.StatusInternalServerError, err)
			}
		}

		next.ServeHTTP(recorder, request)
		duration := time.Since(start).Microseconds()

		_, span := otelTracer.Start(request.Context(), fmt.Sprintf("response body"),
			trace.WithAttributes(
				attribute.String("response.status", strconv.Itoa(recorder.status)),
				attribute.String("response.size", formatSize(recorder.size)),
				attribute.String("response.duration_ms", strconv.FormatInt(duration, 10)),
			))
		if recorder.status == http.StatusOK {
			if recorder.logRespBody {
				span.SetAttributes(
					attribute.String("response.body", recorder.buffer.String()),
				)
			}
		}

		span.End()
	}
}

func queryParamToSpan(r *http.Request, attributes map[string][]string) *http.Request {
	ctx, span := otelTracer.Start(r.Context(), "query parameter")
	r = r.WithContext(ctx)
	defer span.End()

	otelAttributes := make([]attribute.KeyValue, 0, len(attributes))
	for key, values := range attributes {
		for _, value := range values {
			otelAttributes = append(otelAttributes, attribute.String("request.query.params."+key, value))
		}
	}

	span.SetAttributes(otelAttributes...)

	return r
}

func addRequestBodyToSpan(r *http.Request) (*http.Request, error) {
	ctx, span := otelTracer.Start(r.Context(), "request body")
	r = r.WithContext(ctx)
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return r, err
	}
	defer func() {
		errReqBody := r.Body.Close()
		if errReqBody != nil {
			span.RecordError(err)
		}
	}()

	var requestBody map[string]any
	if err = json.Unmarshal(body, &requestBody); err != nil {
		return r, err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	jsonString, err := json.Marshal(requestBody)
	if err != nil {
		return r, err
	}

	span.SetAttributes(attribute.String("request.body.json", string(jsonString)))

	return r, nil
}

func formatSize(size int) string {
	if size < 1024 {
		return fmt.Sprintf("%d bytes", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	}
}
