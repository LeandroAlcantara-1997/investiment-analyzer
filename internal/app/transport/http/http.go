package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LeandroAlcantara-1997/investment-analyzer/config/env"
	docs "github.com/LeandroAlcantara-1997/investment-analyzer/docs"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/app/container"
	analyzer "github.com/LeandroAlcantara-1997/investment-analyzer/internal/app/transport/http/v1/analyzer"
	customcontext "github.com/LeandroAlcantara-1997/investment-analyzer/pkg/custom_context"
	"github.com/LeandroAlcantara-1997/investment-analyzer/pkg/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type api struct {
	allowOrigins string
	port         string
	version      string
	apiName      string
	environment  string
	container    *container.Container
}

func New(port, apiName, version, allowOrigins, environment string, container *container.Container) *api {
	return &api{
		port:         port,
		apiName:      apiName,
		version:      version,
		allowOrigins: allowOrigins,
		environment:  environment,
		container:    container,
	}
}

func (a *api) NewServer(ctx context.Context) {
	r := gin.Default()
	r.ContextWithFallback = true

	r.Use(otelgin.Middleware(env.Env.APIName))

	r.GET("/health-check", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	a.initDoc()
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	r.Use(cors.New(cors.Config{
		AllowOrigins:  util.ChunkTextByComma(a.allowOrigins),
		AllowMethods:  []string{http.MethodGet},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length", "content-type"},
	})).Use(customcontext.AddEndpoint).Use(customcontext.AddHeader).
		Use(customcontext.AddHost).Use(customcontext.AddIP).
		Use(customcontext.AddUserAgent)
	analyzer.ConfigureRoutes(r, a.container.Domains.AnalyzerUseCase)

	log.Printf("Server listening in :%s", a.port)

	a.shutdown(&http.Server{
		Addr:    fmt.Sprintf(":%s", a.port),
		Handler: r,
	})
}

func (a *api) shutdown(server *http.Server) {
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
	ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}

func (a *api) initDoc() {
	docs.SwaggerInfo.Title = a.apiName
	docs.SwaggerInfo.Version = a.version
}
