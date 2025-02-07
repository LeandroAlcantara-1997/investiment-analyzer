package repository

import (
	"context"
	"io"
	"time"

	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/models"
)

//go:generate mockgen -destination ../../mock/repository_mock.go -package=mock -mock_names=Repository=RepositoryMock -source=repository.go
type Repository interface {
	Trader
}

type Trader interface {
	ReadFile(ctx context.Context, key string, file io.ReadCloser) error
	GetPriceCompanyTimeByKey(ctx context.Context, key string, finalDate time.Time) float64
	GetOperations(ctx context.Context, initialDate, finalDate time.Time) []models.Trader
}
