package admin

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestOTELIntegration(t *testing.T) {
	// Set up a test tracer
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}
	
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			t.Errorf("failed to shutdown tracer provider: %v", err)
		}
	}()

	// Test that we can get a tracer
	tracer := GetTracer()
	if tracer == nil {
		t.Error("expected non-nil tracer")
	}

	// Test span creation
	ctx := context.Background()
	ctx, span := StartSpan(ctx, "test.operation")
	defer FinishSpan(span)

	// Test span attributes
	SetSpanAttributes(span, "GET", "/test/endpoint")
	SetSpanResponse(span, 200)

	// The span should be active
	if span == nil {
		t.Error("expected non-nil span")
	}
}

func TestTracerName(t *testing.T) {
	if TracerName != "github.com/ctreminiom/go-atlassian/v2/admin" {
		t.Errorf("expected TracerName to be 'github.com/ctreminiom/go-atlassian/v2/admin', got %s", TracerName)
	}
}