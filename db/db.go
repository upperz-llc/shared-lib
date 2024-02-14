package db

import (
	"context"
	"os/user"
	"time"

	"github.com/upperz-llc/shared-lib/alarm"
	"github.com/upperz-llc/shared-lib/auth"
	"github.com/upperz-llc/shared-lib/device"
	"github.com/upperz-llc/shared-lib/manufacturing"
)

type SQLDB interface {
	// Auth
	AddAuthAndACLs(ctx context.Context, did, username, password string) error
	AddGatewayACLs(ctx context.Context, gid, did string) error
	DeleteGatewayACLs(ctx context.Context, gid, did string) error
	GetACL(ctx context.Context, did, topic string) (*auth.ACL, error)
	GetAuth(ctx context.Context, did string) (*auth.Auth, error)
	// Alarm interface
	// AddNewAlarmToAlarmTimeline(ctx context.Context, alarm Alarm) error
	// CreateAlarmConnection(ctx context.Context, clientID string) (*Alarm, error) // TODO : Too high level
	// IncrementAlarmAckCheckCount(ctx context.Context, alarmID string) error
	// UpdateAlarmAck(ctx context.Context, alarmID, userUID string) error // TODO : Too high level
	// UpdateAlarmTimelineWithClosedAt(ctx context.Context, alarm Alarm) error
	// CloseAlarm(ctx context.Context, alarm *Alarm) error // TODO : Too high level
	// GetDeviceOwner(ctx context.Context, deviceID string) (string, error)

	// GetAuthByDeviceID(ctx context.Context, did string)
	// Alarm interface
	// GetAlarm(ctx context.Context, alarm *Alarm) (*Alarm, error)
	// DeleteAlarm(ctx context.Context, alarmID string) error
	// GetAlarm(ctx context.Context, alarmID string) (*Alarm, error)
	CloseAlarm(ctx context.Context, aid string) (*alarm.Alarm, error)
	CreateAlarm(ctx context.Context, did string, at alarm.AlarmType) (*alarm.Alarm, error)
	QueryAlarm(ctx context.Context, did string, at alarm.AlarmType) (*alarm.Alarm, error)
	QueryAlarmsByUser(ctx context.Context, uid string) ([]alarm.Alarm, error)

	// // Device interface
	// AddDeviceTelemetry(ctx context.Context, deviceID string, data *DeviceTelemetry) error
	// CreateDevice(ctx context.Context, device *Device) error
	// // DeleteDevice(ctx context.Context, deviceID string) error
	CreateDeviceConfig(ctx context.Context, did string, config device.Config) error
	GetDevice(ctx context.Context, did string) (*device.Device, error)
	GetDevicesByOwner(ctx context.Context, uid string) ([]device.Device, error)
	GetDeviceConfig(ctx context.Context, did string) (*device.Config, error)

	UpdateDeviceOwner(ctx context.Context, did, uid string) error
	UpdateDeviceConnectionStatus(ctx context.Context, did string, connectionStatus device.ConnectionStatus) error
	UpdateDeviceFirmwareVersion(ctx context.Context, did, firmwareVersion string) error
	UpdateDeviceMonitoringStatus(ctx context.Context, did string, status device.MonitoringStatus) error
	UpdateDeviceOTAStatus(ctx context.Context, did string, status device.OTAStatus, timestamp int64) error

	// User interface
	CreateUser(ctx context.Context, user user.User) error
	DeleteUserByUID(ctx context.Context, key string) error
	GetUser(ctx context.Context, uid string) (*user.User, error)
	GetUserByEmailAddress(ctx context.Context, email string) (*user.User, error)
	// GetUserNotificationSettings(ctx context.Context, uid string) (*User, error)

	// Device Telemetry
	CreateDeviceTelemetry(ctx context.Context, did string, data device.Telemetry) error
	GetDeviceTelemetry(ctx context.Context, did string, r device.TelemetryRange) ([]device.Telemetry, error)

	// Manufacturing
	CreateDeviceAndManufacturingData(ctx context.Context, md manufacturing.ManufacturingData) error
	CreateManufacturingData(ctx context.Context, md manufacturing.ManufacturingData) error

	// Custom
	GetInactiveGatewayDevices(ctx context.Context, qt time.Time) ([]device.Device, error)

	UpdateLastSeen(ctx context.Context, did string, qt time.Time) error
}
