package sharedlib

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CockroachDB struct {
	pool *pgxpool.Pool
}

func (cdb *CockroachDB) CreateDeviceTelemetry(ctx context.Context, did string, data DeviceTelemetry) error {
	query := `INSERT INTO defaultdb.public.temperature (device_id, temperature, timestamp) VALUES (@device_id, @temperature, @timestamp)`
	args := pgx.NamedArgs{
		"device_id":   did,
		"temperature": math.Round(data.Temperature*100) / 100,
		"timestamp":   time.Unix(data.Timestamp, 0),
	}

	_, err := cdb.pool.Exec(ctx, query, args)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (cbd *CockroachDB) CreateUser(ctx context.Context, user User) error {
	query := `INSERT INTO defaultdb.public.user (uid, email created_at) VALUES (@uid, @email, @created_at)`
	args := pgx.NamedArgs{
		"uid":        user.UID,
		"email":      user.Email,
		"created_at": time.Unix(time.Now().Unix(), 0),
	}

	_, err := cbd.pool.Exec(ctx, query, args)
	if err != nil {
		log.Println(err)
	}

	return err
}

func NewCockroachDB(ctx context.Context) (*CockroachDB, error) {
	dbp := os.Getenv("DB_PASS")
	dsn := "postgresql://temporary:%s@hefty-tiger-10243.5xj.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"

	pool, err := pgxpool.New(ctx, fmt.Sprintf(dsn, dbp))
	if err != nil {
		log.Fatal("failed to create connection pool", err)
		return nil, err
	}

	return &CockroachDB{
		pool: pool,
	}, nil

}
