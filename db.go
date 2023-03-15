package sharedlib

import "context"

type DB interface {
	// Device interface
	GetDeviceOwner(ctx context.Context, deviceID string) (string, error)

	// User interface
	GetUser(ctx context.Context, uid string) (*User, error)
	GetUserNotificationSettings(ctx context.Context, uid string) (*User, error)
}
