package logger

import (
	"context"
)

//go:generate mockgen -destination ../../mock/repository_mock.go -package=mock -mock_names=Repository=RepositoryMock -source=repository.go
type Logger interface {
	Error(ctx context.Context, err error)
	Success(ctx context.Context)
}
