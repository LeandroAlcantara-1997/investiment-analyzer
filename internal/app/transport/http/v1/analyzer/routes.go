package hero

import (
	"net/http"

	"github.com/LeandroAlcantara-1997/investment-analyzer/config/env"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/app/transport/http/middleware"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/output/analyzer"
	"github.com/LeandroAlcantara-1997/investment-analyzer/pkg/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(r *gin.Engine, analyzerService analyzer.Analyzer) {
	analyzer := Handler{
		UseCase: analyzerService,
	}

	m := &middleware.Middleware{
		Admin: false,
		Origin: middleware.Origin{
			Cors: &cors.Config{
				AllowOrigins:  util.ChunkTextByComma(env.Env.AllowOrigins),
				AllowMethods:  []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
				AllowHeaders:  []string{"*"},
				ExposeHeaders: []string{"Content-Length", "content-type"},
			},
		},
	}

	analyzerRoute := r.Group("/v1/reports").Use(m.Init)
	analyzerRoute.GET("", analyzer.getReport)
}
