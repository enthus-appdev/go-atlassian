package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/ctreminiom/go-atlassian/v2/admin"
	"github.com/ctreminiom/go-atlassian/v2/assets"
	"github.com/ctreminiom/go-atlassian/v2/bitbucket"
	"github.com/ctreminiom/go-atlassian/v2/confluence"
	jiragile "github.com/ctreminiom/go-atlassian/v2/jira/agile"
	jirasm "github.com/ctreminiom/go-atlassian/v2/jira/sm"
	jirav2 "github.com/ctreminiom/go-atlassian/v2/jira/v2"
	jirav3 "github.com/ctreminiom/go-atlassian/v2/jira/v3"
	confluencev2 "github.com/ctreminiom/go-atlassian/v2/confluence/v2"
)

// TracingExample demonstrates OTEL integration across all modules
func main() {
	// Set up OpenTelemetry tracer
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
		}
	}()

	ctx := context.Background()

	// Demonstrate tracing across all modules
	fmt.Println("=== OpenTelemetry Integration Test ===")
	
	// Test admin module tracer
	adminTracer := admin.GetTracer()
	ctx, adminSpan := adminTracer.Start(ctx, "test.admin.operation")
	admin.SetSpanAttributes(adminSpan, "GET", "/admin/test")
	admin.SetSpanResponse(adminSpan, 200)
	admin.FinishSpan(adminSpan)
	fmt.Printf("✓ Admin module tracer: %s\n", admin.TracerName)

	// Test assets module tracer
	assetsTracer := assets.GetTracer()
	ctx, assetsSpan := assetsTracer.Start(ctx, "test.assets.operation")
	assets.SetSpanAttributes(assetsSpan, "PUT", "/assets/test")
	assets.SetSpanResponse(assetsSpan, 201)
	assets.FinishSpan(assetsSpan)
	fmt.Printf("✓ Assets module tracer: %s\n", assets.TracerName)

	// Test bitbucket module tracer
	bitbucketTracer := bitbucket.GetTracer()
	ctx, bitbucketSpan := bitbucketTracer.Start(ctx, "test.bitbucket.operation")
	bitbucket.SetSpanAttributes(bitbucketSpan, "POST", "/bitbucket/test")
	bitbucket.SetSpanResponse(bitbucketSpan, 202)
	bitbucket.FinishSpan(bitbucketSpan)
	fmt.Printf("✓ Bitbucket module tracer: %s\n", bitbucket.TracerName)

	// Test confluence module tracer
	confluenceTracer := confluence.GetTracer()
	ctx, confluenceSpan := confluenceTracer.Start(ctx, "test.confluence.operation")
	confluence.SetSpanAttributes(confluenceSpan, "DELETE", "/confluence/test")
	confluence.SetSpanResponse(confluenceSpan, 204)
	confluence.FinishSpan(confluenceSpan)
	fmt.Printf("✓ Confluence module tracer: %s\n", confluence.TracerName)

	// Test confluence v2 module tracer
	confluencev2Tracer := confluencev2.GetTracer()
	ctx, confluencev2Span := confluencev2Tracer.Start(ctx, "test.confluence.v2.operation")
	confluencev2.SetSpanAttributes(confluencev2Span, "PATCH", "/confluence/v2/test")
	confluencev2.SetSpanResponse(confluencev2Span, 200)
	confluencev2.FinishSpan(confluencev2Span)
	fmt.Printf("✓ Confluence v2 module tracer: %s\n", confluencev2.TracerName)

	// Test JIRA module tracers
	jiraAgileTracer := jiragile.GetTracer()
	ctx, jiraAgileSpan := jiraAgileTracer.Start(ctx, "test.jira.agile.operation")
	jiragile.SetSpanAttributes(jiraAgileSpan, "GET", "/jira/agile/test")
	jiragile.SetSpanResponse(jiraAgileSpan, 200)
	jiragile.FinishSpan(jiraAgileSpan)
	fmt.Printf("✓ JIRA Agile module tracer: %s\n", jiragile.TracerName)

	jiraSMTracer := jirasm.GetTracer()
	ctx, jiraSMSpan := jiraSMTracer.Start(ctx, "test.jira.sm.operation")
	jirasm.SetSpanAttributes(jiraSMSpan, "GET", "/jira/sm/test")
	jirasm.SetSpanResponse(jiraSMSpan, 200)
	jirasm.FinishSpan(jiraSMSpan)
	fmt.Printf("✓ JIRA Service Management module tracer: %s\n", jirasm.TracerName)

	jiraV2Tracer := jirav2.GetTracer()
	ctx, jiraV2Span := jiraV2Tracer.Start(ctx, "test.jira.v2.operation")
	jirav2.SetSpanAttributes(jiraV2Span, "GET", "/jira/v2/test")
	jirav2.SetSpanResponse(jiraV2Span, 200)
	jirav2.FinishSpan(jiraV2Span)
	fmt.Printf("✓ JIRA v2 module tracer: %s\n", jirav2.TracerName)

	jiraV3Tracer := jirav3.GetTracer()
	ctx, jiraV3Span := jiraV3Tracer.Start(ctx, "test.jira.v3.operation")
	jirav3.SetSpanAttributes(jiraV3Span, "GET", "/jira/v3/test")
	jirav3.SetSpanResponse(jiraV3Span, 200)
	jirav3.FinishSpan(jiraV3Span)
	fmt.Printf("✓ JIRA v3 module tracer: %s\n", jirav3.TracerName)

	fmt.Println("\n=== All modules successfully configured for OpenTelemetry tracing ===")
}