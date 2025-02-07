package log

import (
	"context"
	"net/http"
	"time"

	customcontext "github.com/LeandroAlcantara-1997/investment-analyzer/pkg/custom_context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

//go:generate mockgen -destination ../../mock/log_mock.go -package=mock -source=log.go
type Logger interface {
	Error(ctx context.Context, err error, data any)
}

type logger struct {
	stdout *zap.Logger
	vendor *pgxpool.Pool
}

func NewLogger(zapLogger *zap.Logger, client *pgxpool.Pool) *logger {
	return &logger{
		stdout: zapLogger,
		vendor: client,
	}
}

type logLevel string

const (
	logLevelError   logLevel = "error"
	logLevelSuccess logLevel = "success"
)

type Log struct {
	ID        string
	DateTime  time.Time
	Host      string
	UserAgent string
	IP        string
	Method    string
	Endpoint  string
	Params    string
	Headers   http.Header
	Message   string
	LogLevel  logLevel
}

func (l *logger) Error(ctx context.Context, err error) {
	data := &Log{
		Host:      customcontext.GetHost(ctx),
		UserAgent: customcontext.GetUserAgent(ctx),
		IP:        customcontext.GetIP(ctx),
		Method:    customcontext.GetMethod(ctx),
		Endpoint:  customcontext.GetEndpoint(ctx),
		Headers:   customcontext.GetHeader(ctx),
		LogLevel:  logLevelError,
		Message:   err.Error(),
	}
	l.stdout.Error(data.Message, zap.Any("data", data))
	if err := l.createLogRegister(ctx, data); err != nil {
		l.stdout.Error(err.Error(), zap.Any("data", data))
	}

}

func (l *logger) Success(ctx context.Context) {
	data := &Log{
		Host:      customcontext.GetHost(ctx),
		UserAgent: customcontext.GetUserAgent(ctx),
		IP:        customcontext.GetIP(ctx),
		Method:    customcontext.GetMethod(ctx),
		Endpoint:  customcontext.GetEndpoint(ctx),
		Headers:   customcontext.GetHeader(ctx),
		LogLevel:  logLevelSuccess,
	}
	l.stdout.Info(data.Message, zap.Any("data", data))
	if err := l.createLogRegister(ctx, data); err != nil {
		l.stdout.Error(err.Error(), zap.Any("data", data))
	}

}

func (l *logger) createLogRegister(ctx context.Context, log *Log) error {
	var (
		query = `INSERT INTO operation (id, consultation_time , ip, query_parameters,
			status)
			VALUES ($1, $2, $3, $4, $5);`
	)

	tx, err := l.vendor.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(err error) error {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				return err
			}
		}
		return tx.Commit(ctx)
	}(err)

	if _, err := tx.Exec(ctx, query, log.ID,
		log.DateTime, log.IP, log.Params, log.LogLevel); err != nil {
		return err
	}
	return nil
}
