package main

import (
	"context"

	"github.com/filhodanuvem/ytgoapi/internal/database"
	"github.com/filhodanuvem/ytgoapi/internal/observability"

	"github.com/filhodanuvem/ytgoapi/internal/http"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	ctx := context.Background()
	otelShutdown, err := observability.SetupOtel(ctx, observability.OtelConfig{
		ServiceName:              "ytgoapi",
		ServiceVersion:           "0.0.1",
		OtelExporterOtlpEndpoint: "otel-collector:4317",
		OtelExporterOtlpInsecure: true,
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		err = otelShutdown(ctx)
		if err != nil {
			panic(err)
		}
	}()

	connectionString := "postgresql://posts:p0stgr3s@db:5432/posts"
	conn, err := database.NewConnection(ctx, connectionString)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(otelgin.Middleware("ytgoapi"))
	g.Use(http.LogMiddleware())
	http.Configure()
	http.SetRoutes(g)
	g.Run(":3000")
}
