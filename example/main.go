package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Lajule/tracing"
	"github.com/Lajule/tracing/example/model"
	"github.com/flamego/flamego"
	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

	db, err := otelsql.Open("mysql", "root:otel_password@tcp(mysql)/db?parseTime=true", otelsql.WithAttributes(semconv.DBSystemMySQL))
	if err != nil {
		panic(err)
	}

	f := flamego.Classic()

	f.Use(flamego.Renderer(
		flamego.RenderOptions{
			JSONIndent: "  ",
		},
	))

	f.Use(func(c flamego.Context) {
		c.Map(model.New(db))
	})

	f.Use(tracing.Tracing())

	f.Get("/", func(c flamego.Context, r flamego.Render, parent context.Context, q *model.Queries) {
		ctx, child := otel.Tracer("handler").Start(parent, "ListAuthors")
		defer child.End()

		authors, err := q.ListAuthors(ctx)
		if err != nil {
			child.SetStatus(codes.Error, "ListAuthors failed")
			child.RecordError(err)

			r.JSON(http.StatusInternalServerError, map[string]string{
				"status": http.StatusText(http.StatusInternalServerError),
			})
			return
		}

		r.JSON(http.StatusOK, authors)
	})

	f.Run()
}
