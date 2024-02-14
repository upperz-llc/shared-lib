package auth

import "time"

type ACL struct {
	ID        string    `json:"id"`
	AuthID    string    `json:"auth_id"`
	DeviceID  string    `json:"device_id"`
	Allowed   bool      `json:"allowed"`
	Topic     string    `json:"topic"`
	Access    string    `json:"Access"`
	CreatedAt time.Time `json:"created_at"`
}

// Auth placeholder
type Auth struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id,omitempty"`
	Enabled   bool      `json:"enabled,omitempty"`
	Username  string    `json:"username,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
