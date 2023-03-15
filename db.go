package sharedlib

import "context"

type DB interface {
	// Alarm interface
	AddNewAlarmToAlarmTimeline(ctx context.Context, alarm Alarm) error
	CreateAlarmConnection(ctx context.Context, clientID string) (*Alarm, error) // TODO : Too high level
	DeleteAlarm(ctx context.Context, alarm *Alarm) error                        // TODO : Too high level
	GetAlarm(ctx context.Context, alarmID string) (*Alarm, error)
	IncrementAlarmAckCheckCount(ctx context.Context, alarmID string) error
	QueryAlarm(ctx context.Context, clientID string, alarmType AlarmType) (Alarm, error)
	UpdateAlarmAck(ctx context.Context, alarmID, userUID string) error // TODO : Too high level
	UpdateAlarmTimelineWithClosedAt(ctx context.Context, alarm Alarm) error

	// Device interface
	GetDeviceOwner(ctx context.Context, deviceID string) (string, error)

	// User interface
	GetUser(ctx context.Context, uid string) (*User, error)
	GetUserNotificationSettings(ctx context.Context, uid string) (*User, error)
}
