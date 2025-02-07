package analyzer

import (
	"context"

	dto "github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/input/analyzer"
)

//go:generate mockgen -destination ../../../mock/hero_mock.go -package=mock -source=service.go
type Analyzer interface {
	GetReport(ctx context.Context, req *dto.AnalyzerRequest) (*dto.AnalyzersResponse, error)
}
