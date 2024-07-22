# tracing

Use [OpenTelemetry][1] with [Flamego][2] framework.


## Installation

	go get github.com/Lajule/tracing

## Getting started

```go
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Lajule/tracing"
	"github.com/flamego/flamego"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	exp, err := otlptrace.New(ctx, otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint("otel-collector:4317"), otlptracegrpc.WithInsecure()))
	if err != nil {
		panic(err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("example"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	defer tp.Shutdown(ctx)

	f := flamego.Classic()
	f.Use(tracing.Tracing())

	f.Get("/", func(parent context.Context) string {
		_, child := otel.Tracer("handler").Start(parent, "Hello")
		defer child.End()
		return "Hello!"
	})

	f.Run()
}
```

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.

[1]: https://opentelemetry.io/
[2]: https://flamego.dev/
