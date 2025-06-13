package assets

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	// TracerName is the name of the tracer for the assets module
	TracerName = "github.com/ctreminiom/go-atlassian/v2/assets"
)

// GetTracer returns an OpenTelemetry tracer for the assets module
func GetTracer() trace.Tracer {
	return otel.Tracer(TracerName)
}

// StartSpan starts a new tracing span with the given name and context
func StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return GetTracer().Start(ctx, spanName, opts...)
}

// SetSpanAttributes sets HTTP request attributes on the span
func SetSpanAttributes(span trace.Span, method, endpoint string) {
	span.SetAttributes(
		attribute.String("http.method", method),
		attribute.String("http.url", endpoint),
		attribute.String("component", "go-atlassian"),
		attribute.String("module", "assets"),
	)
}

// SetSpanError sets error attributes on the span
func SetSpanError(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
}

// SetSpanResponse sets HTTP response attributes on the span
func SetSpanResponse(span trace.Span, statusCode int) {
	span.SetAttributes(
		attribute.Int("http.status_code", statusCode),
	)
	
	if statusCode >= 400 {
		span.SetStatus(codes.Error, http.StatusText(statusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}
}

// FinishSpan properly ends a span
func FinishSpan(span trace.Span) {
	span.End()
}