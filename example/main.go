// Copyright 2022 Flamego. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Lajule/tracing"
	"github.com/Lajule/tracing/example/model"
	"github.com/flamego/flamego"
	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	exp, _ := otlptrace.New(context.Background(), otlptracegrpc.NewClient(otlptracegrpc.WithEndpoint("otel-collector:4317"), otlptracegrpc.WithInsecure()))

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("example"),
		)),
	)

	otel.SetTracerProvider(tp)

	db, _ := otelsql.Open("mysql", "root:otel_password@tcp(mysql)/db?parseTime=true", otelsql.WithAttributes(
		semconv.DBSystemMySQL,
	))

	f := flamego.Classic()

	f.Use(tracing.Tracing(tp.Tracer("middleware")))

	f.Use(func(c flamego.Context) {
		c.Map(model.New(db))
	})

	f.Get("/", func(c flamego.Context, parent context.Context, q *model.Queries) string {
		tracer := otel.Tracer("handler")
		ctx, child := tracer.Start(parent, "Hello")
		defer child.End()

		_, _ = q.ListAuthors(ctx)

		return "Hello, Flamego!"
	})
	f.Run()
}
