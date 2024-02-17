package cockroach

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/upperz-llc/shared-lib/alarm"
	"github.com/upperz-llc/shared-lib/auth"
	"github.com/upperz-llc/shared-lib/device"
	"github.com/upperz-llc/shared-lib/manufacturing"
	"github.com/upperz-llc/shared-lib/user"
)

var once sync.Once
var db *CockroachDB

type CockroachDB struct {
	pool *pgxpool.Pool
}

func (cdb *CockroachDB) CloseAlarm(ctx context.Context, aid string) (*alarm.Alarm, error) {
	query := `UPDATE defaultdb.public.alarm SET active = @active, closed_at = @closed_at WHERE id = @id RETURNING *`
	args := pgx.NamedArgs{
		"id":        aid,
		"active":    false,
		"closed_at": time.Now().Format(time.RFC3339),
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachalarm, err := pgx.CollectOneRow[CockroachAlarm](rows, pgx.RowToStructByPos[CockroachAlarm])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	alarm := cockroachalarm.ToAlarm()
	return &alarm, err

}

func (cdb *CockroachDB) CreateAlarm(ctx context.Context, did string, at alarm.AlarmType) (*alarm.Alarm, error) {
	query := `INSERT INTO defaultdb.public.alarm (id, type, device_id) VALUES (DEFAULT, @type, @device_id) RETURNING *`
	args := pgx.NamedArgs{
		"device_id": did,
		"type":      int(at),
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachalarm, err := pgx.CollectOneRow[CockroachAlarm](rows, pgx.RowToStructByPos[CockroachAlarm])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	alarm := cockroachalarm.ToAlarm()
	return &alarm, err
}

type CockroachAlarm struct {
	ID              pgtype.UUID        `json:"id"`
	Type            pgtype.Int2        `json:"type"`
	DeviceID        pgtype.UUID        `json:"device_id"`
	AckedBy         pgtype.Text        `json:"acked_by"`
	Acked           pgtype.Bool        `json:"acked"`
	Active          pgtype.Bool        `json:"active"`
	AckedCheckCount pgtype.Int2        `json:"acked_check_count"`
	ClosedAt        pgtype.Timestamptz `json:"closed_at"`
	AckedAt         pgtype.Timestamptz `json:"acked_at"`
	CreatedAt       pgtype.Timestamptz `json:"created_at"`
}

func (c CockroachAlarm) ToAlarm() alarm.Alarm {
	d := alarm.Alarm{}
	if c.ID.Valid {
		v, _ := c.ID.Value()
		d.ID = v.(string)
	}
	if c.Type.Valid {
		v, _ := c.Type.Value()
		d.Type = alarm.AlarmType(v.(int64))
	}
	if c.DeviceID.Valid {
		v, _ := c.DeviceID.Value()
		d.DeviceID = v.(string)
	}
	if c.AckedBy.Valid {
		v, _ := c.AckedBy.Value()
		d.AckedBy = v.(string)
	}
	if c.Acked.Valid {
		v, _ := c.Acked.Value()
		d.Acked = v.(bool)
	}
	if c.Active.Valid {
		v, _ := c.Active.Value()
		d.Active = v.(bool)
	}
	if c.AckedCheckCount.Valid {
		v, _ := c.AckedCheckCount.Value()
		d.AckedCheckCount = int(v.(int64))
	}
	if c.ClosedAt.Valid {
		v, _ := c.ClosedAt.Value()
		d.ClosedAt = v.(time.Time)
	}
	if c.AckedAt.Valid {
		v, _ := c.AckedAt.Value()
		d.AckedAt = v.(time.Time)
	}
	if c.CreatedAt.Valid {
		v, _ := c.CreatedAt.Value()
		d.CreatedAt = v.(time.Time)
	}

	return d
}

func (cdb *CockroachDB) QueryAlarm(ctx context.Context, did string, at alarm.AlarmType) (*alarm.Alarm, error) {
	query := `SELECT id, type, device_id, acked_by, acked, active, acked_check_count, closed_at, acked_at, created_at FROM defaultdb.public.alarm WHERE device_id = @device_id and type = @type ORDER BY created_at DESC`
	args := pgx.NamedArgs{
		"device_id": did,
		"type":      int(at),
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachalarm, err := pgx.CollectOneRow[CockroachAlarm](rows, pgx.RowToStructByPos[CockroachAlarm])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	alarm := cockroachalarm.ToAlarm()
	return &alarm, err
}

func (cdb *CockroachDB) QueryAlarmsByUser(ctx context.Context, uid string) ([]alarm.Alarm, error) {
	query := `SELECT id, type, device_id, acked_by, acked, active, acked_check_count, closed_at, acked_at, created_at from  defaultdb.public.alarm where device_id in ( SELECT id from  defaultdb.public.device where owner = @owner) AND active = true;`
	args := pgx.NamedArgs{
		"owner": uid,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachalarms, err := pgx.CollectRows[CockroachAlarm](rows, pgx.RowToStructByPos[CockroachAlarm])

	if errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("no rows")
		return nil, nil
	}

	alarms := make([]alarm.Alarm, 0)
	for _, v := range cockroachalarms {
		alarms = append(alarms, v.ToAlarm())

	}

	return alarms, err
}

// ************ AUTH *******************

type CockroachACL struct {
	ID        pgtype.UUID        `json:"id"`
	AuthID    pgtype.UUID        `json:"auth_id"`
	DeviceID  pgtype.UUID        `json:"device_id"`
	Allowed   pgtype.Bool        `json:"allowed"`
	Topic     pgtype.Text        `json:"topic"`
	Access    pgtype.Text        `json:"Access"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

func (c CockroachACL) ToACL() auth.ACL {
	d := auth.ACL{}
	if c.ID.Valid {
		v, _ := c.ID.Value()
		d.ID = v.(string)
	}
	if c.AuthID.Valid {
		v, _ := c.AuthID.Value()
		d.AuthID = v.(string)
	}
	if c.DeviceID.Valid {
		v, _ := c.DeviceID.Value()
		d.DeviceID = v.(string)
	}
	if c.Allowed.Valid {
		v, _ := c.Allowed.Value()
		d.Allowed = v.(bool)
	}
	if c.Topic.Valid {
		v, _ := c.Topic.Value()
		d.Topic = v.(string)
	}
	if c.Access.Valid {
		v, _ := c.Access.Value()
		d.Access = v.(string)
	}
	if c.CreatedAt.Valid {
		v, _ := c.CreatedAt.Value()
		d.CreatedAt = v.(time.Time)
	}

	return d
}

func (cdb *CockroachDB) GetACL(ctx context.Context, did, topic string) (*auth.ACL, error) {
	query := `SELECT id, auth_id, device_id, allowed, topic, access, created_at FROM defaultdb.public.acl WHERE device_id = @device_id AND topic = @topic`
	args := pgx.NamedArgs{
		"device_id": did,
		"topic":     topic,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachacl, err := pgx.CollectOneRow[CockroachACL](rows, pgx.RowToStructByPos[CockroachACL])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	acl := cockroachacl.ToACL()
	return &acl, nil
}

type CockroachAuth struct {
	ID        pgtype.UUID        `json:"id"`
	DeviceID  pgtype.UUID        `json:"device_id"`
	Enabled   pgtype.Bool        `json:"enabled"`
	Username  pgtype.Text        `json:"username"`
	Password  pgtype.Text        `json:"password"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

func (c CockroachAuth) ToAuth() auth.Auth {
	d := auth.Auth{}
	if c.ID.Valid {
		v, _ := c.ID.Value()
		d.ID = v.(string)
	}
	if c.DeviceID.Valid {
		v, _ := c.DeviceID.Value()
		d.DeviceID = v.(string)
	}
	if c.CreatedAt.Valid {
		v, _ := c.CreatedAt.Value()
		d.CreatedAt = v.(time.Time)
	}
	if c.Enabled.Valid {
		v, _ := c.Enabled.Value()
		d.Enabled = v.(bool)
	}
	if c.Username.Valid {
		v, _ := c.Username.Value()
		d.Username = v.(string)
	}
	if c.Password.Valid {
		v, _ := c.Password.Value()
		d.Password = v.(string)
	}

	return d
}

func (cdb *CockroachDB) GetAuth(ctx context.Context, did string) (*auth.Auth, error) {
	query := `SELECT id, device_id, enabled, username, password, created_at FROM defaultdb.public.auth WHERE device_id = @device_id`
	args := pgx.NamedArgs{
		"device_id": did,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachauth, err := pgx.CollectOneRow[CockroachAuth](rows, pgx.RowToStructByPos[CockroachAuth])

	auth := cockroachauth.ToAuth()
	return &auth, err
}

// func (cdb *CockroachDB) UpdateDeviceOwner(ctx context.Context, did, uid string) error {
// 	query := `UPDATE defaultdb.public.device SET owner = @owner WHERE id = @id`
// 	args := pgx.NamedArgs{
// 		"id":    did,
// 		"owner": uid,
// 	}

// 	_, err := cdb.pool.Exec(ctx, query, args)
// 	return err
// }

// **************************************

func (cdb *CockroachDB) CreateDeviceConfig(ctx context.Context, did string, config device.Config) error {
	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	query := `DELETE FROM defaultdb.public.device_config WHERE device_id = @device_id`
	args := pgx.NamedArgs{
		"device_id": did,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	query = `INSERT INTO defaultdb.public.device_config (id, device_id, alert, warning, target, measurement_interval, version, created_at, updated_at) VALUES (DEFAULT, @device_id, @alert, @warning, @target, @measurement_interval, @version, DEFAULT, NOW())`
	args = pgx.NamedArgs{
		"device_id":            did,
		"alert":                config.AlertTemperature,
		"warning":              config.WarningTemperature,
		"target":               config.TargetTemperature,
		"measurement_interval": config.TelemetryPeriod,
		"version":              config.Version,
	}

	_, err = cdb.pool.Exec(ctx, query, args)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

// ************* DEVICE *****************
type CockroachDevice struct {
	ID               pgtype.UUID `json:"id"`
	ConnectionStatus pgtype.Int8 `json:"connection_status"`
	Type             pgtype.Int8 `json:"type"`
	FirmwareVersion  pgtype.Text `json:"firmware_version"`
	MonitoringStatus pgtype.Int8 `json:"monitoring_status"`
	Owner            pgtype.UUID `json:"owner"`
}

func (c CockroachDevice) ToDevice() device.Device {
	d := device.Device{}
	if c.ID.Valid {
		v, _ := c.ID.Value()
		d.ID = v.(string)
	}
	if c.ConnectionStatus.Valid {
		v, _ := c.ConnectionStatus.Value()
		d.ConnectionStatus = device.ConnectionStatus(v.(int64))
	}
	if c.Type.Valid {
		v, _ := c.Type.Value()
		d.Type = device.Type(v.(int64))
	}
	if c.FirmwareVersion.Valid {
		v, _ := c.FirmwareVersion.Value()
		d.FirmwareVersion = v.(string)
	}
	if c.MonitoringStatus.Valid {
		v, _ := c.MonitoringStatus.Value()
		d.MonitoringStatus = device.MonitoringStatus(v.(int64))
	}
	if c.Owner.Valid {
		v, _ := c.Owner.Value()
		d.Owner = v.(string)
	}
	return d
}

func (cdb *CockroachDB) GetDevice(ctx context.Context, did string) (*device.Device, error) {
	query := `SELECT id, connection_status, device_type, firmware_version, monitoring_status, owner FROM defaultdb.public.device WHERE id = @device_id`
	args := pgx.NamedArgs{
		"device_id": did,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachdevice, err := pgx.CollectExactlyOneRow[CockroachDevice](rows, pgx.RowToStructByPos[CockroachDevice])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	device := cockroachdevice.ToDevice()
	return &device, err
}

type CockroachDeviceConfig struct {
	ID                  pgtype.UUID
	DeviceID            pgtype.UUID
	Alert               pgtype.Int8
	Warning             pgtype.Int8
	Target              pgtype.Int8
	MeasurementInterval pgtype.Int8
	Version             pgtype.Int8
	CreatedAt           pgtype.Timestamptz
	Updated_at          pgtype.Timestamptz
}

func (c CockroachDeviceConfig) ToDeviceConfig() device.Config {
	d := device.Config{}
	if c.Alert.Valid {
		v, _ := c.Alert.Value()
		d.AlertTemperature = int(v.(int64))
	}
	if c.Warning.Valid {
		v, _ := c.Warning.Value()
		d.WarningTemperature = int(v.(int64))
	}
	if c.Target.Valid {
		v, _ := c.Target.Value()
		d.TargetTemperature = int(v.(int64))
	}
	if c.MeasurementInterval.Valid {
		v, _ := c.MeasurementInterval.Value()
		d.TelemetryPeriod = int(v.(int64))
	}
	if c.Version.Valid {
		v, _ := c.Version.Value()
		d.Version = int(v.(int64))
	}

	return d
}

func (cdb *CockroachDB) GetDeviceConfig(ctx context.Context, did string) (*device.Config, error) {
	query := `SELECT id, device_id, alert, warning, target, measurement_interval, version, created_at, updated_at FROM defaultdb.public.device_config WHERE device_id = @device_id`
	args := pgx.NamedArgs{
		"device_id": did,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachdeviceconfig, err := pgx.CollectOneRow[CockroachDeviceConfig](rows, pgx.RowToStructByPos[CockroachDeviceConfig])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	deviceconfig := cockroachdeviceconfig.ToDeviceConfig()
	return &deviceconfig, err
}

func (cdb *CockroachDB) GetDevicesByOwner(ctx context.Context, uid string) ([]device.Device, error) {
	query := `SELECT id, connection_status, device_type, firmware_version, monitoring_status, nickname, temperature, owner, last_seen FROM defaultdb.public.device WHERE owner = @uid`
	args := pgx.NamedArgs{
		"uid": uid,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachdevices, err := pgx.CollectRows[CockroachDevice](rows, pgx.RowToStructByPos[CockroachDevice])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	devices := make([]device.Device, 0)
	for _, v := range cockroachdevices {
		devices = append(devices, v.ToDevice())

	}

	return devices, err
}

func (cdb *CockroachDB) UpdateDeviceOwner(ctx context.Context, did, uid string) error {
	query := `UPDATE defaultdb.public.device SET owner = @owner WHERE id = @id`
	args := pgx.NamedArgs{
		"id":    did,
		"owner": uid,
	}

	_, err := cdb.pool.Exec(ctx, query, args)
	return err
}

// ************************************

func (cdb *CockroachDB) UpdateDeviceFirmwareVersion(ctx context.Context, did, firmwareVersion string) error {
	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {

		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}

	query := `UPDATE defaultdb.public.device SET last_seen = @timestamp, firmware_version = @firmware_version WHERE id = @id`
	args := pgx.NamedArgs{
		"id":               did,
		"firmware_version": firmwareVersion,
		"timestamp":        time.Now().Format(time.RFC3339),
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

func (cdb *CockroachDB) UpdateDeviceConnectionStatus(ctx context.Context, did string, status device.ConnectionStatus) error {
	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	query := `UPDATE defaultdb.public.device SET updated_at = @timestamp, connection_status = @connection_status WHERE id = @id`
	args := pgx.NamedArgs{
		"id":                did,
		"timestamp":         time.Now().Format(time.RFC3339),
		"connection_status": status,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

func (cdb *CockroachDB) UpdateDeviceMonitoringStatus(ctx context.Context, did string, status device.MonitoringStatus) error {
	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {

		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}

	query := `UPDATE defaultdb.public.device SET last_seen = @timestamp, monitoring_status = @monitoring_status WHERE id = @id`
	args := pgx.NamedArgs{
		"id":                did,
		"timestamp":         time.Now().Format(time.RFC3339),
		"monitoring_status": status,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

func (cdb *CockroachDB) UpdateDeviceOTAStatus(ctx context.Context, did string, status device.OTAStatus, timestamp int64) error {
	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {

		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}

	query := `INSERT INTO defaultdb.public.ota_status (id, device_id, created_at, status) VALUES (DEFAULT, @device_id, DEFAULT, @status)`
	args := pgx.NamedArgs{
		"device_id": did,
		"status":    status,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	query = `UPDATE defaultdb.public.device SET last_seen = @timestamp WHERE id = @id`
	args = pgx.NamedArgs{
		"id":        did,
		"timestamp": time.Now(),
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

// func (cdb *CockroachDB) CreateDeviceTelemetry(ctx context.Context, did string, data device.Telemetry) error {
// 	conn, err := cdb.pool.Acquire(ctx)
// 	if err != nil {

// 		return err
// 	}
// 	defer conn.Release()

// 	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
// 	if err != nil {

// 		return err
// 	}

// 	query := `INSERT INTO defaultdb.public.temperature (device_id, temperature, timestamp) VALUES (@device_id, @temperature, @timestamp)`
// 	args := pgx.NamedArgs{
// 		"device_id":   did,
// 		"temperature": data.Temperature,
// 		"timestamp":   data.Timestamp,
// 	}

// 	_, err = tx.Exec(ctx, query, args)
// 	if err != nil {
// 		if err := tx.Rollback(ctx); err != nil {
// 			return err
// 		}
// 		return err
// 	}

// 	query = `UPDATE defaultdb.public.device SET last_seen = @timestamp, temperature = @temperature WHERE id = @id`
// 	args = pgx.NamedArgs{
// 		"id":          did,
// 		"temperature": data.Temperature,
// 		"timestamp":   data.Timestamp,
// 	}

// 	_, err = tx.Exec(ctx, query, args)
// 	if err != nil {
// 		if err := tx.Rollback(ctx); err != nil {
// 			return err
// 		}
// 		return err
// 	}

// 	return tx.Commit(ctx)
// }

func (cdb *CockroachDB) CreateUser(ctx context.Context, user user.User) error {
	query := `INSERT INTO defaultdb.public.user (uid, email, password, created_at) VALUES (@uid, @email, @password, @created_at)`
	args := pgx.NamedArgs{
		"uid":        user.UID,
		"email":      user.Email,
		"password":   user.Password,
		"created_at": time.Unix(time.Now().Unix(), 0).Format(time.RFC3339),
	}

	_, err := cdb.pool.Exec(ctx, query, args)

	return err
}

func (cdb *CockroachDB) DeleteUserByUID(ctx context.Context, uid string) error {
	query := `DELETE FROM defaultdb.public.user WHERE uid = @uid`
	args := pgx.NamedArgs{
		"uid": uid,
	}

	_, err := cdb.pool.Exec(ctx, query, args)

	return err
}

// asdadasdadaddadadasdas

func (cdb *CockroachDB) CreateManufacturingData(ctx context.Context, md manufacturing.ManufacturingData) error {
	query := `INSERT INTO defaultdb.public.device_manufacturing_info (device_id, device_type, manufactured_at, measurement_type, username, password) VALUES (@device_id, @device_type, @manufactured_at, @measurement_type, @username, @password)`
	args := pgx.NamedArgs{
		"device_id":        md.DeviceID,
		"device_type":      md.DeviceType,
		"manufactured_at":  time.Unix(time.Now().Unix(), 0).Format(time.RFC3339),
		"measurement_type": md.MeasurementType,
		"username":         md.Username,
		"password":         md.Password,
	}

	_, err := cdb.pool.Exec(ctx, query, args)

	return err
}

func (cdb *CockroachDB) CreateDevice(ctx context.Context, deviceID, username, password string, deviceType device.Type, measurementType device.MeasurementType) error {
	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	query := `INSERT INTO defaultdb.public.auth (id, enabled, username, password, device_id) VALUES (DEFAULT, @enabled, @username, @password, @device_id) RETURNING id`
	args := pgx.NamedArgs{
		"enabled":   true,
		"username":  username,
		"password":  password,
		"device_id": deviceID,
	}
	var aid pgtype.UUID
	if err := tx.QueryRow(ctx, query, args).Scan(&aid); err != nil {
		return err
	}

	query = `INSERT INTO defaultdb.public.device (id, device_type, measurement_type, connection_status) VALUES (@device_id, @device_type, @measurement_type, @connection_status)`
	args = pgx.NamedArgs{
		"device_id":         deviceID,
		"device_type":       deviceType,
		"measurement_type":  measurementType,
		"connection_status": device.Disconnected,
	}
	if _, err := tx.Exec(ctx, query, args); err != nil {
		return err
	}

	query = `INSERT INTO defaultdb.public.acl (id, auth_id, device_id, topic, access, allowed) VALUES
	(DEFAULT, @auth_id, @device_id, @topic1, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic2, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic3, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic4, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic5, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic6, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic7, @access, @allowed)`
	args = pgx.NamedArgs{
		"auth_id":   aid,
		"device_id": deviceID,
		"topic1":    fmt.Sprintf("DATA/%s", deviceID),
		"topic2":    fmt.Sprintf("CMD/%s", deviceID),
		"topic3":    fmt.Sprintf("CMD/%s/response", deviceID),
		"topic4":    fmt.Sprintf("CONFIG/%s", deviceID),
		"topic5":    fmt.Sprintf("CONFIG/%s/response", deviceID),
		"topic6":    fmt.Sprintf("STATE/%s", deviceID),
		"topic7":    fmt.Sprintf("LWT/%s", deviceID),
		"access":    "rw",
		"allowed":   true,
	}

	// query = `INSERT INTO defaultdb.public.device_manufacturing_data (id, device_id, device_type, manufactured_at, measurement_type, username, password) VALUES (DEFAULT, @device_id, @device_type, @manufactured_at, @measurement_type, @username, @password) RETURNING id`
	// args = pgx.NamedArgs{
	// 	"device_id":        md.DeviceID,
	// 	"device_type":      md.DeviceType,
	// 	"measurement_type": md.MeasurementType,
	// 	"manufactured_at":  md.ManufacturedAt,
	// 	"username":         md.Username,
	// 	"password":         md.Password,
	// }
	if _, err := tx.Exec(ctx, query, args); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

func (cdb *CockroachDB) AddAuthAndACLs(ctx context.Context, did, username, password string) error {
	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {

		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}

	query := `INSERT INTO defaultdb.public.auth (id, enabled, username, password) VALUES (DEFAULT, @enabled, @username, @password) RETURNING id`
	args := pgx.NamedArgs{
		"enabled":  true,
		"username": username,
		"password": password,
	}
	var aid pgtype.UUID
	if err := tx.QueryRow(ctx, query, args).Scan(&aid); err != nil {
		return err
	}

	query = `INSERT INTO defaultdb.public.acl (id, auth_id, device_id, topic, access, allowed) VALUES
	(DEFAULT, @auth_id, @device_id, @topic1, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic2, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic3, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic4, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic5, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic6, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic7, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic8, @access, @allowed)`
	args = pgx.NamedArgs{
		"auth_id":   aid,
		"device_id": did,
		"topic1":    fmt.Sprintf("DATA/%s", did),
		"topic2":    fmt.Sprintf("CMD/%s", did),
		"topic3":    fmt.Sprintf("BCMD/%s", did),
		"topic4":    fmt.Sprintf("BCMD/%s/response", did),
		"topic5":    fmt.Sprintf("CONFIG/%s", did),
		"topic6":    fmt.Sprintf("CONFIG/%s/response", did),
		"topic7":    fmt.Sprintf("STATE/%s", did),
		"topic8":    fmt.Sprintf("LWT/%s", did),
		"access":    "rw",
		"allowed":   true,
	}
	if _, err := tx.Exec(ctx, query, args); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

func (cdb *CockroachDB) AddGatewayACLs(ctx context.Context, gid, did string) error {

	auth, err := cdb.GetAuth(ctx, gid)
	if err != nil {

		return err
	}

	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {

		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}

	query := `INSERT INTO defaultdb.public.acl (id, auth_id, device_id, topic, access, allowed) VALUES
	(DEFAULT, @auth_id, @device_id, @topic1, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic2, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic3, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic4, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic5, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic6, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic7, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic8, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic8, @access, @allowed),
	(DEFAULT, @auth_id, @device_id, @topic9, @access, @allowed)`
	args := pgx.NamedArgs{
		"auth_id":   auth.ID,
		"device_id": gid,
		"topic1":    fmt.Sprintf("DATA/%s", did),
		"topic2":    fmt.Sprintf("CMD/%s", did),
		"topic3":    fmt.Sprintf("BCMD/%s", did),
		"topic4":    fmt.Sprintf("BCMD/%s/response", did),
		"topic5":    fmt.Sprintf("CONFIG/%s", did),
		"topic6":    fmt.Sprintf("CONFIG/%s/response", did),
		"topic7":    fmt.Sprintf("STATE/%s", did),
		"topic8":    fmt.Sprintf("BIRTH/%s", did),
		"topic9":    fmt.Sprintf("DEATH/%s", did),
		"topic10":   fmt.Sprintf("LWT/%s", did),
		"access":    "rw",
		"allowed":   true,
	}
	if _, err := tx.Exec(ctx, query, args); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

func (cdb *CockroachDB) DeleteGatewayACLs(ctx context.Context, gid, did string) error {

	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {

		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}

	query := `DELETE FROM defaultdb.public.acl WHERE device_id = @gateway_id AND topic LIKE @device_id`
	args := pgx.NamedArgs{
		"gateway_id": gid,
		"device_id":  fmt.Sprintf("%%%s%%", did),
	}
	if _, err := tx.Exec(ctx, query, args); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

type CockroachUser struct {
	ID               pgtype.UUID        `json:"id"`
	UID              pgtype.Text        `json:"uid"`
	Email            pgtype.Text        `json:"email"`
	Password         pgtype.Text        `json:"password"`
	NotificationPush pgtype.Bool        `json:"notification_push"`
	NotificationSMS  pgtype.Bool        `json:"notification_sms"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	PhoneNumber      pgtype.Text        `json:"phone_number"`
}

func (c CockroachUser) ToUser() user.User {
	d := user.User{}
	if c.ID.Valid {
		v, _ := c.ID.Value()
		d.ID = v.(string)
	}
	if c.UID.Valid {
		v, _ := c.UID.Value()
		d.UID = v.(string)
	}
	if c.Email.Valid {
		v, _ := c.Email.Value()
		d.Email = v.(string)
	}
	if c.Password.Valid {
		v, _ := c.Password.Value()
		d.Password = v.(string)
	}
	if c.NotificationPush.Valid {
		v, _ := c.NotificationPush.Value()
		d.NotificationPush = v.(bool)
	}
	if c.NotificationSMS.Valid {
		v, _ := c.NotificationSMS.Value()
		d.NotificationSMS = v.(bool)
	}
	if c.CreatedAt.Valid {
		v, _ := c.CreatedAt.Value()
		d.CreatedAt = v.(time.Time)
	}
	if c.UpdatedAt.Valid {
		v, _ := c.UpdatedAt.Value()
		d.UpdatedAt = v.(time.Time)
	}
	if c.PhoneNumber.Valid {
		v, _ := c.PhoneNumber.Value()
		d.PhoneNumber = v.(string)
	}

	return d
}

func (cdb *CockroachDB) GetUser(ctx context.Context, uid string) (*user.User, error) {
	query := `SELECT id, uid, email, password, notification_push, notification_sms, created_at, updated_at, phone_number FROM defaultdb.public.user WHERE uid = @uid`
	args := pgx.NamedArgs{
		"uid": uid,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachuser, err := pgx.CollectOneRow[CockroachUser](rows, pgx.RowToStructByPos[CockroachUser])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	user := cockroachuser.ToUser()
	return &user, err
}

func (cdb *CockroachDB) GetUserByEmailAddress(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, uid, email, password, notification_push, notification_sms, created_at, updated_at, phone_number FROM defaultdb.public.user WHERE email = @email`
	args := pgx.NamedArgs{
		"email": email,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachuser, err := pgx.CollectOneRow[CockroachUser](rows, pgx.RowToStructByPos[CockroachUser])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	user := cockroachuser.ToUser()
	return &user, err
}

type CockroachTelemetry struct {
	ID          pgtype.UUID        `json:"id"`
	DeviceID    pgtype.UUID        `json:"device_id"`
	Temperature pgtype.Int8        `json:"temperature"`
	Timestamp   pgtype.Timestamptz `json:"timestamp"`
}

// func (c CockroachTelemetry) ToDeviceTelemetry() device.Telemetry {
// 	d := device.Temperature{}
// 	if c.ID.Valid {
// 		v, _ := c.ID.Value()
// 		d.ID = v.(string)
// 	}
// 	if c.DeviceID.Valid {
// 		v, _ := c.DeviceID.Value()
// 		d.DeviceID = v.(string)
// 	}
// 	if c.Temperature.Valid {
// 		v, _ := c.Temperature.Value()
// 		d.Temperature = v.(int64)
// 	}
// 	if c.Timestamp.Valid {
// 		v, _ := c.Timestamp.Value()
// 		d.Timestamp = v.(time.Time)
// 	}

// 	return d
// }

// func (cdb *CockroachDB) GetDeviceTelemetry(ctx context.Context, did string, r device.TelemetryRange) ([]device.Telemetry, error) {
// 	query := `SELECT id, device_id, temperature, timestamp FROM defaultdb.public.temperature WHERE device_id = @device_id AND timestamp >= @range`
// 	args := pgx.NamedArgs{
// 		"device_id": did,
// 		"range":     r.ToTime(),
// 	}

// 	rows, err := cdb.pool.Query(ctx, query, args)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cockroachtelemetry, err := pgx.CollectRows[CockroachTelemetry](rows, pgx.RowToStructByPos[CockroachTelemetry])

// 	if errors.Is(err, pgx.ErrNoRows) {
// 		return nil, nil
// 	}

// 	telemetry := make([]device.Telemetry, 0)
// 	for _, v := range cockroachtelemetry {
// 		telemetry = append(telemetry, v.ToDeviceTelemetry())
// 	}

// 	return telemetry, err
// }

func (cdb *CockroachDB) GetInactiveGatewayDevices(ctx context.Context, qt time.Time) ([]device.Device, error) {
	query := `SELECT id, connection_status, device_type, firmware_version, monitoring_status, nickname, temperature, owner, last_seen FROM defaultdb.public.device WHERE connection_status = 'connected' AND device_type = '2'
	AND last_seen < @time`
	args := pgx.NamedArgs{
		"time": qt,
	}

	rows, err := cdb.pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	cockroachdevices, err := pgx.CollectRows[CockroachDevice](rows, pgx.RowToStructByPos[CockroachDevice])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	devices := make([]device.Device, 0)
	for _, v := range cockroachdevices {
		devices = append(devices, v.ToDevice())

	}

	return devices, err
}

func (cdb *CockroachDB) UpdateLastSeen(ctx context.Context, did string, timestamp time.Time) error {
	conn, err := cdb.pool.Acquire(ctx)
	if err != nil {

		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}

	query := `UPDATE defaultdb.public.device SET last_seen = @timestamp WHERE id = @id`
	args := pgx.NamedArgs{
		"id":        did,
		"timestamp": time.Now(),
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return err
		}
		return err
	}

	return tx.Commit(ctx)
}

func NewCockroachDB(ctx context.Context) (*CockroachDB, error) {

	once.Do(func() {
		dbu := os.Getenv("DB_USERNAME")
		dbp := os.Getenv("DB_PASS")
		dbc := os.Getenv("DB_CLUSTER")
		// dsn := "postgresql://manufacturing:%s@hefty-tiger-10243.5xj.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
		dsct := "postgresql://%s:%s@%s/defaultdb?sslmode=verify-full"
		dscs := fmt.Sprintf(dsct, dbu, dbp, dbc)

		config, err := pgxpool.ParseConfig(dscs)
		if err != nil {
			panic(err)
		}

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			panic(err)
		}

		db = &CockroachDB{
			pool: pool,
		}
	})

	return db, nil
}
