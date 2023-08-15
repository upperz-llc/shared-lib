package sharedlib

import (
	"context"

	"github.com/upperz-llc/shared-lib/alarm"
)

// type DB interface {
// 	// Alarm interface
// 	AddNewAlarmToAlarmTimeline(ctx context.Context, alarm Alarm) error
// 	CreateAlarmConnection(ctx context.Context, clientID string) (*Alarm, error) // TODO : Too high level
// 	IncrementAlarmAckCheckCount(ctx context.Context, alarmID string) error
// 	UpdateAlarmAck(ctx context.Context, alarmID, userUID string) error // TODO : Too high level
// 	UpdateAlarmTimelineWithClosedAt(ctx context.Context, alarm Alarm) error
// 	CloseAlarm(ctx context.Context, alarm *Alarm) error // TODO : Too high level
// 	GetDeviceOwner(ctx context.Context, deviceID string) (string, error)

// 	// Alarm interface
// 	CreateAlarm(ctx context.Context, alarm *Alarm) error
// 	DeleteAlarm(ctx context.Context, alarmID string) error
// 	GetAlarm(ctx context.Context, alarmID string) (*Alarm, error)
// 	QueryAlarm(ctx context.Context, clientID string, alarmType AlarmType) (*Alarm, error)

// 	// Device interface
// 	AddDeviceTelemetry(ctx context.Context, deviceID string, data *DeviceTelemetry) error
// 	CreateDevice(ctx context.Context, device *Device) error
// 	// DeleteDevice(ctx context.Context, deviceID string) error
// 	GetDevice(ctx context.Context, deviceID string) (*Device, error)
// 	UpdateDeviceConnectionStatus(ctx context.Context, deviceID string, connectionStatus DeviceConnectionStatus) error
// 	UpdateDeviceFirmwareVersion(ctx context.Context, deviceID, firmwareVersion string) error
// 	UpdateDeviceOTAStatus(ctx context.Context, deviceID string, status OTAStatus, timestamp int64) error

// 	// User interface
// 	CreateUser(ctx context.Context, user *User) error
// 	DeleteUser(ctx context.Context, uid string) error
// 	GetUser(ctx context.Context, uid string) (*User, error)
// 	GetUserNotificationSettings(ctx context.Context, uid string) (*User, error)
// }

type SQLDB interface {
	// Auth
	AddAuthAndACLs(ctx context.Context, did, username, password string) error
	GetACL(ctx context.Context, did, topic string) (*ACL, error)
	GetAuth(ctx context.Context, did string) (*Auth, error)
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
	CreateAlarm(ctx context.Context, did string, at alarm.AlarmType) error
	QueryAlarm(ctx context.Context, did string, at alarm.AlarmType) (*alarm.Alarm, error)

	// // Device interface
	// AddDeviceTelemetry(ctx context.Context, deviceID string, data *DeviceTelemetry) error
	// CreateDevice(ctx context.Context, device *Device) error
	// // DeleteDevice(ctx context.Context, deviceID string) error
	CreateDeviceConfig(ctx context.Context, did string, config DeviceConfig) error
	GetDevice(ctx context.Context, did string) (*Device, error)
	GetDevicesByOwner(ctx context.Context, uid string) ([]Device, error)
	GetDeviceConfig(ctx context.Context, did string) (*DeviceConfig, error)
	UpdateDeviceOwner(ctx context.Context, did, uid string) error
	UpdateDeviceConnectionStatus(ctx context.Context, did string, connectionStatus DeviceConnectionStatus) error
	UpdateDeviceFirmwareVersion(ctx context.Context, did, firmwareVersion string) error
	UpdateDeviceMonitoringStatus(ctx context.Context, did string, status DeviceMonitoringStatus) error
	UpdateDeviceOTAStatus(ctx context.Context, did string, status OTAStatus, timestamp int64) error

	// User interface
	CreateUser(ctx context.Context, user User) error
	DeleteUserByUID(ctx context.Context, key string) error
	GetUser(ctx context.Context, uid string) (*User, error)
	// GetUserNotificationSettings(ctx context.Context, uid string) (*User, error)

	// Device Telemetry
	CreateDeviceTelemetry(ctx context.Context, did string, data DeviceTelemetry) error
	GetDeviceTelemetry(ctx context.Context, did string, r TelemetryRange) ([]DeviceTelemetry, error)

	// Manufacturing
	CreateDeviceAndManufacturingData(ctx context.Context, md ManufacturingData) error
	CreateManufacturingData(ctx context.Context, md ManufacturingData) error
}
