package hero

import (
	"net/http"

	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/app/transport/http/response"
	dto "github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/input/analyzer"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/output/analyzer"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

type Handler struct {
	UseCase analyzer.Analyzer
}

// @Summary      Get a investiment report
// @Description  Get a investiment report
// @Tags         Analyzers
// @Accept       json
// @Produce      json
// @Param hero body dto.AnalyzerRequest true "analyzer"
// @Success      200  {object}  dto.AnalyzersResponse
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /reports [get]
func (h *Handler) getReport(ctx *gin.Context) {
	c, span := otel.Tracer("analyzer").Start(ctx.Request.Context(), "getReport")
	defer span.End()
	var request dto.AnalyzerRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.UseCase.GetReport(c, &request)
	if err != nil {
		ctx.JSON(response.RestError(err))
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}
