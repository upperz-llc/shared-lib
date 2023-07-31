package sharedlib

// Auth placeholder
type Auth struct {
	ID        string `json:"id"`
	DeviceID  string `json:"device_id,omitempty"`
	Enabled   bool   `json:"enabled,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
