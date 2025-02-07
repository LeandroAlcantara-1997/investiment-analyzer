package repository

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/exception"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/models"
)

var (
	operations     = make([]models.Trader, 0)
	companiesPrice = make(map[string][]models.CompanyPrice, 0)
)

func (r *repo) GetOperations(ctx context.Context, initialDate, finalDate time.Time) []models.Trader {
	operation := make([]models.Trader, 0)
	for _, p := range operations {
		if p.ActionTime.After(initialDate) &&
			p.ActionTime.Before(finalDate) {
			operation = append(operations, p)
		}
	}
	return operation
}

func (r *repo) GetPriceCompanyTimeByKey(ctx context.Context, key string, finalDate time.Time) float64 {
	operationCompany := make([]models.CompanyPrice, 0)
	for _, c := range companiesPrice[key] {
		if c.PriceTime.Before(finalDate) || c.PriceTime.Equal(finalDate) {
			operationCompany = append(operationCompany, c)
		}
	}
	return operationCompany[len(operationCompany)-1].Price
}

func (r *repo) ReadFile(ctx context.Context, key string, file io.ReadCloser) error {
	var initalLine = true

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if initalLine {
			initalLine = false
			continue
		}

		if err := createTrader(ctx,
			strings.Split(string(line), ","), key); err != nil {
			return exception.New(fmt.Sprintf("createTrader -> %s", err.Error()), err)
		}

	}

	return nil
}

type processor interface {
	processFile(ctx context.Context, operation []string) (err error)
}

func createTrader(ctx context.Context, operation []string, key string) (err error) {
	var p processor
	switch len(operation) {
	case 5:
		p = new(actionsByTime)
	case 2:
		p = &priceCompanyByTime{key: key}
	default:
		return exception.ErrInvalidFile
	}

	return p.processFile(ctx, operation)
}

type priceCompanyByTime struct {
	key string
}

func (p *priceCompanyByTime) processFile(ctx context.Context, operation []string) (err error) {
	var company = new(models.CompanyPrice)
	company.PriceTime, err = time.Parse(time.DateTime, operation[0])
	if err != nil {
		return err
	}

	company.Price, err = strconv.ParseFloat(operation[1], 64)
	if err != nil {
		return exception.New(fmt.Sprintf("parseFloat -> %s", err.Error()), err)
	}

	companiesPrice[p.key] = append(companiesPrice[p.key], *company)
	return
}

type actionsByTime struct {
}

func (a *actionsByTime) processFile(ctx context.Context, operation []string) (err error) {
	var trader = new(models.Trader)
	trader.ActionTime, err = time.Parse(time.DateTime, operation[0])
	if err != nil {
		return err
	}

	trader.CompanyName = operation[1]

	trader.Quantity, err = strconv.Atoi(operation[2])
	if err != nil {
		return exception.New(fmt.Sprintf("atoi -> %s", err.Error()), err)
	}
	trader.Price, err = strconv.ParseFloat(operation[3], 64)
	if err != nil {
		return exception.New(fmt.Sprintf("parseFloat -> %s", err.Error()), err)
	}
	trader.Side = operation[4]
	operations = append(operations, *trader)
	return
}
