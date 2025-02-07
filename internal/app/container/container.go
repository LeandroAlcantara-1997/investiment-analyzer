package container

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"github.com/LeandroAlcantara-1997/investment-analyzer/config/env"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/adapter/cache"
	log "github.com/LeandroAlcantara-1997/investment-analyzer/internal/adapter/log"
	"github.com/LeandroAlcantara-1997/investment-analyzer/internal/adapter/repository"
	analyzerimplementation "github.com/LeandroAlcantara-1997/investment-analyzer/internal/domain/analyzer/service"
	analyzer "github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/output/analyzer"
	repositorycontract "github.com/LeandroAlcantara-1997/investment-analyzer/internal/ports/output/repository"
	"github.com/LeandroAlcantara-1997/investment-analyzer/pkg/otel"
	"github.com/exaring/otelpgx"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Container struct {
	Domains    Domains
	components *components
}

type Domains struct {
	AnalyzerUseCase analyzer.Analyzer
}

type components struct {
	pgxClient   *pgxpool.Pool
	redisClient *redis.Client
	zapLogger   *zap.Logger
	repository  repositorycontract.Repository
}

func New() (context.Context, *Container, error) {
	env.LoadEnv()
	ctx := context.Background()
	otel.New(env.Env.APIName, env.Env.Environment).TraceProvider(ctx)

	cmp, err := setupComponents(ctx)
	if err != nil {
		return ctx, nil, fmt.Errorf("setupComponents -> %w", err)
	}
	analyzerService := analyzerimplementation.New(
		cmp.repository,
		cache.New(cmp.redisClient),
		log.NewLogger(cmp.zapLogger, cmp.pgxClient),
		env.Env.CashInHand,
	)
	c := &Container{
		Domains: Domains{
			AnalyzerUseCase: analyzerService,
		},
		components: cmp,
	}

	if err := c.readFiles(ctx); err != nil {
		return ctx, nil, err
	}

	return ctx, c, nil
}

func setupComponents(ctx context.Context) (*components, error) {

	pgxClient, err := createConnectionDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("createConnectionDatabase -> %w", err)
	}

	redisClient, err := createConnectionRedis(ctx)
	if err != nil {
		return nil, fmt.Errorf("createConnectionRedis -> %w", err)
	}
	return &components{
		pgxClient:   pgxClient,
		redisClient: redisClient,
		repository:  repository.New(),
	}, nil
}

func (c *Container) GetZapLogger() *zap.Logger {
	return c.components.zapLogger
}

func createConnectionDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	var conn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		env.Env.DBUser,
		env.Env.DBPassword,
		env.Env.DBHost,
		env.Env.DBPort,
		env.Env.DBName,
	)
	pgxConfig, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, fmt.Errorf("parseConfig -> %w", err)
	}
	pgxConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pgxClient, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("newWithConfig -> %w", err)
	}
	if err = pgxClient.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping -> %w", err)
	}
	sqlDB, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("open -> %w", err)
	}
	databaseDriver, err := pgx.WithInstance(sqlDB, &pgx.Config{
		DatabaseName: env.Env.DBName,
	})
	if err != nil {
		return nil, fmt.Errorf("withInstance -> %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./config/migration",
		env.Env.DBName,
		databaseDriver,
	)
	if err != nil {
		return nil, fmt.Errorf("newWithDatabaseInstance -> %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("up -> %w", err)
	}
	return pgxClient, nil
}

func createConnectionRedis(ctx context.Context) (*redis.Client, error) {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", env.Env.CacheHost, env.Env.CachePort),
			Password:     env.Env.CachePassword,
			DB:           0,
			ReadTimeout:  time.Duration(env.Env.CacheReadTimeout) * time.Second,
			WriteTimeout: time.Duration(env.Env.CacheWriteTimeout) * time.Second,
		},
	)
	redisotel.InstrumentTracing(redisClient)
	redisotel.InstrumentMetrics(redisClient)
	cmd := redisClient.Ping(ctx)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("ping -> %w", cmd.Err())
	}
	return redisClient, nil
}

func (c *Container) readFiles(ctx context.Context) error {
	folderPath := "./config/data"
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		f, err := os.ReadFile(folderPath + "/" + file.Name())
		if err != nil {
			return err
		}
		if err := c.components.repository.ReadFile(ctx, string(regexp.MustCompile(`[^A-Z]+`).
			ReplaceAll([]byte(file.Name()), []byte(``))),
			io.NopCloser(bytes.NewReader(f))); err != nil {
			return err
		}
	}

	return nil
}
