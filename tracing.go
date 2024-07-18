package tracing

import (
	"github.com/flamego/flamego"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func Tracing(tracer trace.Tracer) flamego.Handler {
	return flamego.ContextInvoker(func(c flamego.Context) {
		r := c.Request()

		options := []trace.SpanStartOption{
			trace.WithAttributes(attribute.KeyValue{Key: attribute.Key("method"), Value: attribute.StringValue(r.Method)}),
			trace.WithAttributes(attribute.KeyValue{Key: attribute.Key("path"), Value: attribute.StringValue(r.URL.Path)}),
		}

		var remoteSpanCtx trace.SpanContext

		if c.Request().Header != nil {
			propagator := propagation.TraceContext{}
			remoteSpanCtx = trace.SpanContextFromContext(propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header)))
		}

		operation := r.Method + " " + r.URL.Path

		var span trace.Span

		if remoteSpanCtx.IsValid() {
			_, span = tracer.Start(trace.ContextWithRemoteSpanContext(r.Context(), remoteSpanCtx), operation, options...)
		} else {
			_, span = tracer.Start(r.Context(), operation, options...)
		}

		defer span.End()

		c.Map(span)
		c.Next()
	})
}
