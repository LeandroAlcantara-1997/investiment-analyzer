package cache

import "context"

//go:generate mockgen -destination ../../mock/cache_mock.go -package=mock -source=cache.go
type Cache interface {
	GetCompanyID(ctx context.Context, name string) (string, error)
	SetCompanyID(ctx context.Context, name, id string) (err error)
}
