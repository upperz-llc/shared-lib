package user

import "time"

type User struct {
	ID               string    `json:"id"`
	UID              string    `json:"uid"`
	Email            string    `json:"email,omitempty"`
	PhoneNumber      string    `json:"phone_number,omitempty"`
	Password         string    `json:"password,omitempty"`
	NotificationSMS  bool      `json:"notification_sms,omitempty"`
	NotificationPush bool      `json:"notification_push,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}
