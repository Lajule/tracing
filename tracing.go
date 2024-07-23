package tracing

import (
	"github.com/flamego/flamego"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Tracing() flamego.Handler {
	return flamego.ContextInvoker(func(c flamego.Context) {
		r := c.Request()

		ctx, span := otel.Tracer("middleware").Start(
			r.Context(),
			r.Method+" "+r.URL.Path,
			trace.WithAttributes(attribute.KeyValue{Key: attribute.Key("method"), Value: attribute.StringValue(r.Method)}),
			trace.WithAttributes(attribute.KeyValue{Key: attribute.Key("path"), Value: attribute.StringValue(r.URL.Path)}),
			trace.WithAttributes(attribute.KeyValue{Key: attribute.Key("remote-addr"), Value: attribute.StringValue(c.RemoteAddr())}),
		)
		defer span.End()

		c.Map(ctx)
		c.Map(span)
		c.Next()
	})
}
