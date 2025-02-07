package analyzer

import (
	"context"

	"time"

	"go.opentelemetry.io/otel"

	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/adapter/cache"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/models"
	dto "github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/input/analyzer"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/output/logger"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/output/repository"
)

type service struct {
	repository repository.Repository
	cache      cache.Cache
	logger     logger.Logger
	cashInHand float64
}

func New(repository repository.Repository, cache cache.Cache,
	logger logger.Logger, cashInHand float64) *service {
	return &service{
		repository: repository,
		cache:      cache,
		logger:     logger,
		cashInHand: cashInHand,
	}
}
func (s *service) GetReport(ctx context.Context, req *dto.AnalyzerRequest) (*dto.AnalyzersResponse, error) {
	ctx, span := otel.Tracer("analyzer").Start(ctx, "getReport")
	defer span.End()

	resp, err := s.getReport(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *service) getReport(ctx context.Context, req *dto.AnalyzerRequest) (*dto.AnalyzersResponse, error) {
	traders := s.repository.GetOperations(ctx, req.InitialDate.Time, req.FinalDate.Time)

	registers := s.calculateHeritageBySide(ctx, req, traders)
	calculateAccumulatedProfitability(registers)

	return registers, nil
}

func calculateAccumulatedProfitability(registers *dto.AnalyzersResponse) {
	for i := range registers.Registers {
		if i > 0 && i < len(registers.Registers) {
			registers.Registers[i].AccumulatedProfitability = dto.Money(registers.Registers[i-1].HeritageEvolution/
				registers.Registers[i].HeritageEvolution - 1).Round(5)
		}
	}
}

func (s *service) calculateHeritageBySide(ctx context.Context, req *dto.AnalyzerRequest, traders []models.Trader) *dto.AnalyzersResponse {
	var (
		heritage  = s.cashInHand
		dateLimit = req.InitialDate.Add(time.Duration(req.Interval) * time.Minute)
		register  = new(dto.AnalyzersResponse)
		actions   = make(map[string]int, 2)
		companies = make(map[string]float64, 0)
	)

	register.Registers = make([]dto.Register, 0)
	register.Registers = append(register.Registers, dto.Register{
		HeritageEvolution: dto.Money(heritage).Round(2),
		Timestamp:         req.InitialDate,
	})

	for _, trader := range traders {

		if _, ok := companies[trader.CompanyName]; !ok {
			companies[trader.CompanyName] = s.repository.
				GetPriceCompanyTimeByKey(ctx, trader.CompanyName, dateLimit)
		}

		if trader.ActionTime.After(dateLimit) {
			for k, v := range actions {
				heritage += float64(v) * companies[k]
			}
			companies = map[string]float64{}
			dateLimit = dateLimit.Add(time.Minute *
				time.Duration(req.Interval))

			register.Registers = append(register.Registers, dto.Register{
				HeritageEvolution: dto.Money(heritage).Round(2),
				Timestamp:         &dto.Date{Time: dateLimit},
			})
		}

		if trader.Side == "BUY" {
			actions[trader.CompanyName] += trader.Quantity

			heritage -= trader.Price * float64(trader.Quantity)
			continue
		}

		actions[trader.CompanyName] -= trader.Quantity
		heritage += trader.Price
	}

	return register
}
