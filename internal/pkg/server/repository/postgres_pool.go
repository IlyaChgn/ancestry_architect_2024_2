package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresPool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) (err error)
	Close()
}

const (
	defaultMaxConns          = int32(90)
	defaultMinConns          = int32(0)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

func NewConnectionString(user, password, host, port, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
}

func postgresPoolConfig(dbURL string) *pgxpool.Config {
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(_ context.Context, _ *pgx.Conn) bool {
		return true
	}

	dbConfig.AfterRelease = func(_ *pgx.Conn) bool {
		return true
	}

	return dbConfig
}

func NewPostgresPool(dbURL string) (*pgxpool.Pool, error) {
	postgresCfg := postgresPoolConfig(dbURL)

	pool, err := pgxpool.ConnectConfig(context.Background(), postgresCfg)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
