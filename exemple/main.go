// Copyright 2022 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	"github.com/Lajule/tracing"
	"github.com/flamego/flamego"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	exp, _ := otlptrace.New(context.Background(), otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure()))

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("FlamegoService"),
		)),
	)

	otel.SetTracerProvider(tp)

	f := flamego.Classic()
	f.Use(tracing.Tracing(tp.Tracer("Flamego")))
	f.Get("/", func() string {
		return "Hello, Flamego!"
	})
	f.Run()
}
