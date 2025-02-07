package repository

import "github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/output/repository"

type repo struct {
}

func New() repository.Repository {
	return &repo{}
}
