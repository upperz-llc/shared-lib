package sharedlib

import "context"

type DB interface {
	// Alarm interface
	AddNewAlarmToAlarmTimeline(ctx context.Context, alarm Alarm) error
	CreateAlarmConnection(ctx context.Context, clientID string) (*Alarm, error) // TODO : Too high level
	IncrementAlarmAckCheckCount(ctx context.Context, alarmID string) error
	UpdateAlarmAck(ctx context.Context, alarmID, userUID string) error // TODO : Too high level
	UpdateAlarmTimelineWithClosedAt(ctx context.Context, alarm Alarm) error
	CloseAlarm(ctx context.Context, alarm *Alarm) error // TODO : Too high level
	GetDeviceOwner(ctx context.Context, deviceID string) (string, error)

	// Alarm interface
	// CreateAlarm(ctx context.Context, alarm *Alarm) error
	DeleteAlarm(ctx context.Context, alarmID string) error
	GetAlarm(ctx context.Context, alarmID string) (*Alarm, error)
	QueryAlarm(ctx context.Context, clientID string, alarmType AlarmType) (*Alarm, error)

	// Device interface
	AddDeviceTelemetry(ctx context.Context, deviceID string, data *DeviceTelemetry) error
	CreateDevice(ctx context.Context, device *Device) error
	DeleteDevice(ctx context.Context, deviceID string) error
	GetDevice(ctx context.Context, deviceID string) (*Device, error)
	UpdateDeviceConnectionStatus(ctx context.Context, deviceID string, connectionStatus DeviceConnectionStatus) error
	UpdateDeviceFirmwareVersion(ctx context.Context, deviceID, firmwareVersion string) error
	UpdateDeviceOTAStatus(ctx context.Context, deviceID, status string, timestamp int64) error

	// User interface
	// CreateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, uid string) error
	GetUser(ctx context.Context, uid string) (*User, error)
	GetUserNotificationSettings(ctx context.Context, uid string) (*User, error)
}
