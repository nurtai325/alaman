package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurtai325/alaman/internal/config"
)

var (
	pool *pgxpool.Pool
)

const (
	driver = "postgres"
)

func connect(conf config.Config) error {
	dbUrl := fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		driver,
		conf.POSTGRES_USER,
		conf.POSTGRES_PASSWORD,
		conf.POSTGRES_HOST,
		conf.POSTGRES_PORT,
		conf.POSTGRES_DB,
	)
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return err
	}
	pool = conn
	return nil
}

func New(conf config.Config) (*pgxpool.Pool, error) {
	if pool == nil {
		err := connect(conf)
		if err != nil {
			return nil, err
		}
	}
	return pool, nil
}
