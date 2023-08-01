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
	"golang.org/x/exp/slog"
)

type CockroachDB struct {
	pool   *pgxpool.Pool
	logger slog.Logger
}

// ************* DEVICE *****************

func (cdb *CockroachDB) GetDevice(ctx context.Context, did string) (*Device, error) {
	query := `SELECT * FROM defaultdb.public.device WHERE id = @device_id`
	args := pgx.NamedArgs{
		"device_id": did,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	device, err := pgx.CollectOneRow[Device](rows, pgx.RowToStructByName[Device])

	return &device, err
}

func (cdb *CockroachDB) UpdateDeviceOwner(ctx context.Context, did, uid string) error {
	query := `UPDATE defaultdb.public.device (id, owner) VALUES (@id, @owner)`
	args := pgx.NamedArgs{
		"id":    did,
		"owner": uid,
	}

	_, err := cdb.pool.Exec(ctx, query, args)
	return err
}

// ************************************

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

func (cdb *CockroachDB) CreateUser(ctx context.Context, user User) error {
	query := `INSERT INTO defaultdb.public.user (uid, email, created_at) VALUES (@uid, @email, @created_at)`
	args := pgx.NamedArgs{
		"uid":        user.UID,
		"email":      user.Email,
		"created_at": time.Unix(time.Now().Unix(), 0),
	}

	_, err := cdb.pool.Exec(ctx, query, args)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (cdb *CockroachDB) DeleteUserByUID(ctx context.Context, uid string) error {
	query := `DELETE FROM defaultdb.public.user WHERE uid = @uid`
	args := pgx.NamedArgs{
		"uid": uid,
	}

	_, err := cdb.pool.Exec(ctx, query, args)
	if err != nil {
		log.Println(err)
	}

	return err
}

// asdadasdadaddadadasdas

func (cdb *CockroachDB) CreateManufacturingData(ctx context.Context, md ManufacturingData) error {
	query := `INSERT INTO defaultdb.public.device_manufacturing_info (device_id, device_type, manufactured_at, measurement_type, username, password) VALUES (@device_id, @device_type, @manufactured_at, @measurement_type, @username, @password)`
	args := pgx.NamedArgs{
		"device_id":        md.DeviceID,
		"device_type":      md.DeviceType,
		"manufactured_at":  time.Unix(time.Now().Unix(), 0),
		"measurement_type": md.MeasurementType,
		"username":         md.Username,
		"password":         md.Password,
	}

	_, err := cdb.pool.Exec(ctx, query, args)
	if err != nil {
		log.Println(err)
	}

	return err
}

func NewCockroachDB(ctx context.Context) (*CockroachDB, error) {
	dbu := os.Getenv("DB_USERNAME")
	dbp := os.Getenv("DB_PASS")
	dbc := os.Getenv("DB_CLUSTER")
	// dsn := "postgresql://manufacturing:%s@hefty-tiger-10243.5xj.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	dsct := "postgresql://%s:%s@%s/defaultdb?sslmode=verify-full"
	dscs := fmt.Sprintf(dsct, dbu, dbp, dbc)

	config, err := pgxpool.ParseConfig(dscs)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal("failed to create connection pool", err)
		return nil, err
	}

	return &CockroachDB{
		pool: pool,
	}, nil

}
